import { render, screen, waitFor } from "@testing-library/react";
import { AdminDashboard } from "@/components/admin-dashboard";
import * as adminController from "@/controllers/adminController";

jest.mock("@/controllers/adminController", () => ({
  handleLoadAdminDashboard: jest.fn(),
}));

describe("AdminDashboard", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("shows loading state initially", () => {
    (adminController.handleLoadAdminDashboard as jest.Mock).mockImplementation(
      () => new Promise(() => {}) // never resolves
    );

    render(<AdminDashboard />);

    expect(screen.getByText("Loading dashboard...")).toBeInTheDocument();
  });

  it("loads dashboard message and renders dashboard content", async () => {
    (adminController.handleLoadAdminDashboard as jest.Mock).mockResolvedValue({
      ok: true,
      data: {
        message: "Welcome admin!",
      },
    });

    render(<AdminDashboard />);

    expect(await screen.findByText("Admin Dashboard")).toBeInTheDocument();
    expect(screen.getByText("Welcome admin!")).toBeInTheDocument();
  });

  it("renders dashboard even when API fails", async () => {
    (adminController.handleLoadAdminDashboard as jest.Mock).mockResolvedValue({
      ok: false,
    });

    render(<AdminDashboard />);

    await waitFor(() => {
      expect(screen.getByText("Admin Dashboard")).toBeInTheDocument();
    });

    // message should remain empty but dashboard still loads
    expect(screen.getByText("Recent Activity")).toBeInTheDocument();
  });

  it("calls admin controller on mount", async () => {
    (adminController.handleLoadAdminDashboard as jest.Mock).mockResolvedValue({
      ok: true,
      data: { message: "Hello" },
    });

    render(<AdminDashboard />);

    await waitFor(() => {
      expect(
        adminController.handleLoadAdminDashboard
      ).toHaveBeenCalledTimes(1);
    });
  });
});
