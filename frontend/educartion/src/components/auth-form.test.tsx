import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

import { AuthForm } from "@/components/auth-form";
import { handleLogin, handleRegister } from "@/controllers/authController";

const pushMock = jest.fn();

jest.mock("next/navigation", () => ({
  useRouter: () => ({
    push: pushMock,
  }),
}));

jest.mock("@/controllers/authController", () => ({
  handleLogin: jest.fn(),
  handleRegister: jest.fn(),
}));

describe("AuthForm", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe("Login Mode", () => {
    it("renders login form", () => {
      render(<AuthForm mode="login" />);

      expect(screen.getByText("Sign in")).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "Login" })).toBeInTheDocument();
      expect(screen.getByText("Remember me on this device")).toBeInTheDocument();

      expect(
        screen.queryByPlaceholderText("Jane Doe")
      ).not.toBeInTheDocument();
    });

    it("shows error when login fails", async () => {
      (handleLogin as jest.Mock).mockResolvedValue({
        ok: false,
        error: "Invalid credentials",
      });

      const user = userEvent.setup();

      render(<AuthForm mode="login" />);

      await user.type(
        screen.getByPlaceholderText("you@example.com"),
        "test@test.com"
      );

      await user.type(
        screen.getByPlaceholderText("••••••••"),
        "password"
      );

      await user.click(
        screen.getByRole("button", { name: "Login" })
      );

      expect(
        await screen.findByText("Invalid credentials")
      ).toBeInTheDocument();

      expect(pushMock).not.toHaveBeenCalled();
    });

    it("redirects normal user to products after successful login", async () => {
      (handleLogin as jest.Mock).mockResolvedValue({
        ok: true,
        data: {
          full_name: "John Smith",
          role: "Customer",
        },
      });

      const user = userEvent.setup();

      render(<AuthForm mode="login" />);

      await user.type(
        screen.getByPlaceholderText("you@example.com"),
        "john@test.com"
      );

      await user.type(
        screen.getByPlaceholderText("••••••••"),
        "password"
      );

      await user.click(
        screen.getByRole("button", { name: "Login" })
      );

      await waitFor(() => {
        expect(pushMock).toHaveBeenCalledWith("/products");
      });

      expect(handleLogin).toHaveBeenCalledWith({
        email: "john@test.com",
        password: "password",
        rememberMe: false,
      });
    });

    it("redirects admin user to admin page", async () => {
      (handleLogin as jest.Mock).mockResolvedValue({
        ok: true,
        data: {
          full_name: "Admin User",
          role: "Admin",
        },
      });

      const user = userEvent.setup();

      render(<AuthForm mode="login" />);

      await user.type(
        screen.getByPlaceholderText("you@example.com"),
        "admin@test.com"
      );

      await user.type(
        screen.getByPlaceholderText("••••••••"),
        "password"
      );

      await user.click(
        screen.getByRole("button", { name: "Login" })
      );

      await waitFor(() => {
        expect(pushMock).toHaveBeenCalledWith("/admin");
      });
    });
  });

  describe("Register Mode", () => {
    it("renders register form", () => {
      render(<AuthForm mode="register" />);

      expect(
        screen.getByRole("button", { name: "Register" })
      ).toBeInTheDocument();

      expect(
        screen.getByPlaceholderText("Jane Doe")
      ).toBeInTheDocument();

      expect(
        screen.getByText("Confirm password")
      ).toBeInTheDocument();
    });

    it("shows validation error when passwords do not match", async () => {
      const user = userEvent.setup();

      render(<AuthForm mode="register" />);

      await user.type(
        screen.getByPlaceholderText("Jane Doe"),
        "John Smith"
      );

      await user.type(
        screen.getByPlaceholderText("you@example.com"),
        "john@test.com"
      );

      const passwordInputs =
        screen.getAllByPlaceholderText("••••••••");

      await user.type(passwordInputs[0], "password1");
      await user.type(passwordInputs[1], "password2");

      await user.click(
        screen.getByRole("button", { name: "Register" })
      );

      expect(
        await screen.findByText("Passwords do not match.")
      ).toBeInTheDocument();

      expect(handleRegister).not.toHaveBeenCalled();
    });

    it("calls register API with correct payload", async () => {
      (handleRegister as jest.Mock).mockResolvedValue({
        ok: true,
      });

      const user = userEvent.setup();

      render(<AuthForm mode="register" />);

      await user.type(
        screen.getByPlaceholderText("Jane Doe"),
        "John Smith"
      );

      await user.type(
        screen.getByPlaceholderText("you@example.com"),
        "john@test.com"
      );

      const passwordInputs =
        screen.getAllByPlaceholderText("••••••••");

      await user.type(passwordInputs[0], "password");
      await user.type(passwordInputs[1], "password");

      await user.click(
        screen.getByRole("button", { name: "Register" })
      );

      await waitFor(() => {
        expect(handleRegister).toHaveBeenCalledWith({
          fullName: "John Smith",
          email: "john@test.com",
          password: "password",
        });
      });

      expect(pushMock).toHaveBeenCalledWith("/login");
    });

    it("shows registration error", async () => {
      (handleRegister as jest.Mock).mockResolvedValue({
        ok: false,
        error: "Email already exists",
      });

      const user = userEvent.setup();

      render(<AuthForm mode="register" />);

      await user.type(
        screen.getByPlaceholderText("Jane Doe"),
        "John Smith"
      );

      await user.type(
        screen.getByPlaceholderText("you@example.com"),
        "john@test.com"
      );

      const passwordInputs =
        screen.getAllByPlaceholderText("••••••••");

      await user.type(passwordInputs[0], "password");
      await user.type(passwordInputs[1], "password");

      await user.click(
        screen.getByRole("button", { name: "Register" })
      );

      expect(
        await screen.findByText("Email already exists")
      ).toBeInTheDocument();

      expect(pushMock).not.toHaveBeenCalled();
    });
  });
});
