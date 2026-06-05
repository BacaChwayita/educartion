import { login, register } from "@/services/authService";
import { readBackendError } from "@/lib/auth-errors";

jest.mock("@/lib/auth-errors", () => ({
  readBackendError: jest.fn(),
}));

describe("authService", () => {
  beforeEach(() => {
    jest.clearAllMocks();
    global.fetch = jest.fn();
  });

  describe("login", () => {
    it("should return user data on successful login", async () => {
      (fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: async () => ({
          id: 1,
          full_name: "John Doe",
          role: "user",
        }),
      });

      const result = await login({
        email: "test@test.com",
        password: "123456",
        rememberMe: false,
      } as any);

      expect(result.ok).toBe(true);
      expect(result.data?.full_name).toBe("John Doe");
      expect(fetch).toHaveBeenCalledWith("/api/auth/login", expect.any(Object));
    });

    it("should return backend error on failed login", async () => {
      (readBackendError as jest.Mock).mockResolvedValue("Invalid credentials");

      (fetch as jest.Mock).mockResolvedValue({
        ok: false,
        status: 401,
      });

      const result = await login({
        email: "test@test.com",
        password: "wrong",
      } as any);

      expect(result.ok).toBe(false);
      expect(result.error).toBe("Invalid credentials");
      expect(result.status).toBe(401);
    });

    it("should handle network error", async () => {
      (fetch as jest.Mock).mockRejectedValue(new Error("Network fail"));

      const result = await login({
        email: "test@test.com",
        password: "123456",
      } as any);

      expect(result.ok).toBe(false);
      expect(result.error).toBe(
        "Network error while contacting auth endpoint."
      );
    });
  });

  describe("register", () => {
    it("should succeed on valid registration", async () => {
      (fetch as jest.Mock).mockResolvedValue({
        ok: true,
      });

      const result = await register({
        fullName: "John Doe",
        email: "test@test.com",
        password: "123456",
      } as any);

      expect(result.ok).toBe(true);
      expect(fetch).toHaveBeenCalledWith("/api/auth/register", expect.any(Object));
    });

    it("should return backend error on failed registration", async () => {
      (readBackendError as jest.Mock).mockResolvedValue("Email already exists");

      (fetch as jest.Mock).mockResolvedValue({
        ok: false,
        status: 409,
      });

      const result = await register({
        fullName: "John Doe",
        email: "test@test.com",
        password: "123456",
      } as any);

      expect(result.ok).toBe(false);
      expect(result.error).toBe("Email already exists");
      expect(result.status).toBe(409);
    });

    it("should handle network failure", async () => {
      (fetch as jest.Mock).mockRejectedValue(new Error("Network down"));

      const result = await register({
        fullName: "John Doe",
        email: "test@test.com",
        password: "123456",
      } as any);

      expect(result.ok).toBe(false);
      expect(result.error).toBe(
        "Network error while contacting register endpoint."
      );
    });
  });
});
