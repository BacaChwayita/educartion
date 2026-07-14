import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import FeaturedProducts from "@/components/FeaturedProducts";

const products = [
  {
    id: 1,
    title: "Laptop",
    price: 999.99,
    image: "/laptop.jpg",
    description: "Gaming laptop",
  },
  {
    id: 2,
    title: "Mouse",
    price: 49.99,
    image: "/mouse.jpg",
  },
];

describe("FeaturedProducts", () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it("renders all products", () => {
    render(<FeaturedProducts products={products} />);

    expect(screen.getByText("Laptop")).toBeInTheDocument();
    expect(screen.getByText("Mouse")).toBeInTheDocument();

    expect(screen.getByText("$999.99")).toBeInTheDocument();
    expect(screen.getByText("$49.99")).toBeInTheDocument();

    expect(screen.getByText("Gaming laptop")).toBeInTheDocument();
  });

  it("uses product image alt text", () => {
    render(<FeaturedProducts products={products} />);

    expect(screen.getByAltText("Laptop")).toBeInTheDocument();
    expect(screen.getByAltText("Mouse")).toBeInTheDocument();
  });

  it("calls onAddToCart when add to cart button is clicked", async () => {
    const user = userEvent.setup();
    const onAddToCart = jest.fn();

    render(
      <FeaturedProducts
        products={products}
        onAddToCart={onAddToCart}
      />
    );

    const buttons = screen.getAllByRole("button", {
      name: "Add to cart",
    });

    await user.click(buttons[0]);

    expect(onAddToCart).toHaveBeenCalledTimes(1);
    expect(onAddToCart).toHaveBeenCalledWith(products[0]);
  });

  it("does not fail when onAddToCart is not provided", async () => {
    const user = userEvent.setup();

    render(<FeaturedProducts products={products} />);

    const buttons = screen.getAllByRole("button", {
      name: "Add to cart",
    });

    await user.click(buttons[0]);

    expect(screen.getByText("Laptop")).toBeInTheDocument();
  });

  it("scrolls left when left arrow is clicked", async () => {
    const user = userEvent.setup();

    const scrollByMock = jest.fn();

    Object.defineProperty(HTMLElement.prototype, "scrollBy", {
      configurable: true,
      value: scrollByMock,
    });

    render(<FeaturedProducts products={products} />);

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

    render(<FeaturedProducts products={products} />);

    const buttons = screen.getAllByRole("button");

    await user.click(buttons[buttons.length - 1]);

    expect(scrollByMock).toHaveBeenCalledWith({
      left: 400,
      behavior: "smooth",
    });
  });
});
