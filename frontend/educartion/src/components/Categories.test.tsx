import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Categories from "@/components/Categories";

const categories = [
  {
    id: 1,
    name: "Electronics",
    image: "/electronics.jpg",
  },
  {
    id: 2,
    name: "Books",
  },
];

describe("Categories", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders all category names", () => {
    render(<Categories categories={categories} />);

    expect(screen.getByText("Electronics")).toBeInTheDocument();
    expect(screen.getByText("Books")).toBeInTheDocument();
  });

  it("renders category image when provided", () => {
    render(<Categories categories={categories} />);

    expect(
      screen.getByAltText("Electronics")
    ).toBeInTheDocument();
  });

  it("calls onSelect when category is clicked", async () => {
    const user = userEvent.setup();
    const onSelect = jest.fn();

    render(
      <Categories
        categories={categories}
        onSelect={onSelect}
      />
    );

    await user.click(
      screen.getByRole("button", { name: /electronics/i })
    );

    expect(onSelect).toHaveBeenCalledTimes(1);
    expect(onSelect).toHaveBeenCalledWith(1);
  });

  it("does not fail when onSelect is not provided", async () => {
    const user = userEvent.setup();

    render(<Categories categories={categories} />);

    await user.click(
      screen.getByRole("button", { name: /electronics/i })
    );

    expect(
      screen.getByText("Electronics")
    ).toBeInTheDocument();
  });

  it("scrolls left when left arrow is clicked", async () => {
    const user = userEvent.setup();

    const scrollByMock = jest.fn();

    Object.defineProperty(HTMLElement.prototype, "scrollBy", {
      configurable: true,
      value: scrollByMock,
    });

    render(<Categories categories={categories} />);

    const buttons = screen.getAllByRole("button");

    await user.click(buttons[0]);

    expect(scrollByMock).toHaveBeenCalledWith({
      left: -400,
      behavior: "smooth",
    });
  });

  it("scrolls right when right arrow is clicked", async () => {
    const user = userEvent.setup();

    const scrollByMock = jest.fn();

    Object.defineProperty(HTMLElement.prototype, "scrollBy", {
      configurable: true,
      value: scrollByMock,
    });

    render(<Categories categories={categories} />);

    const buttons = screen.getAllByRole("button");

    await user.click(buttons[buttons.length - 1]);

    expect(scrollByMock).toHaveBeenCalledWith({
      left: 400,
      behavior: "smooth",
    });
  });

  it("renders placeholder when category image is missing", () => {
    render(<Categories categories={categories} />);

    expect(screen.queryByAltText("Books")).not.toBeInTheDocument();
    expect(screen.getByText("Books")).toBeInTheDocument();
  });
});
