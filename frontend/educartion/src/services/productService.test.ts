import { getProducts } from "@/services/productService";

describe("getProducts", () => {
  beforeEach(() => {
    global.fetch = jest.fn();
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  it("returns products when API call succeeds", async () => {
    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => ({
        products: [
          { id: 1, name: "Laptop" },
          { id: 2, name: "Mouse" },
        ],
      }),
    });

    const res = await getProducts();

    expect(res.ok).toBe(true);
    expect(res.data).toHaveLength(2);
  });

  it("returns error when API fails", async () => {
    (global.fetch as jest.Mock).mockResolvedValue({
      ok: false,
      status: 500,
      json: async () => ({ error: "Server error" }),
    });

    const res = await getProducts();

    expect(res.ok).toBe(false);
    expect(res.error).toBe("Server error");
    expect(res.status).toBe(500);
  });

  it("handles network failure", async () => {
    (global.fetch as jest.Mock).mockRejectedValue(new Error("Network down"));

    const res = await getProducts();

    expect(res.ok).toBe(false);
    expect(res.error).toContain("Network down");
  });

  it("normalizes array response", async () => {
    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: async () => [
        { id: 1, name: "Laptop" },
        { id: 2, name: "Mouse" },
      ],
    });

    const res = await getProducts();

    expect(res.data).toHaveLength(2);
  });
});
