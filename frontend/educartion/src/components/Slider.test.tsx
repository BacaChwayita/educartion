import { render, screen, act } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import Slider from "@/components/Slider";

const slides = [
  {
    id: 1,
    image: "/img1.jpg",
    title: "Slide 1",
    subtitle: "First slide",
  },
  {
    id: 2,
    image: "/img2.jpg",
    title: "Slide 2",
    subtitle: "Second slide",
  },
];

describe("Slider", () => {
  afterEach(() => {
    jest.useRealTimers();
    jest.clearAllMocks();
  });

  it("returns nothing when no slides are provided", () => {
    const { container } = render(<Slider slides={[]} />);
    expect(container).toBeEmptyDOMElement();
  });

  it("renders first slide initially", () => {
    render(<Slider slides={slides} />);

    expect(screen.getByAltText("Slide 1")).toBeInTheDocument();
    expect(screen.getByText("Slide 1")).toBeInTheDocument();
    expect(screen.getByText("First slide")).toBeInTheDocument();
  });

  it("shows next slide when next button is clicked", async () => {
    const user = userEvent.setup();

    render(<Slider slides={slides} />);

    await user.click(screen.getByRole("button", { name: "Next" }));

    expect(screen.getByAltText("Slide 2")).toBeInTheDocument();
    expect(screen.getByText("Slide 2")).toBeInTheDocument();
  });

  it("shows previous slide when prev button is clicked", async () => {
    const user = userEvent.setup();

    render(<Slider slides={slides} />);

    await user.click(screen.getByRole("button", { name: "Previous" }));

    expect(screen.getByAltText("Slide 2")).toBeInTheDocument();
  });

  it("auto-advances slides using interval", () => {
    jest.useFakeTimers();

    render(<Slider slides={slides} interval={1000} />);

    expect(screen.getByAltText("Slide 1")).toBeInTheDocument();

    act(() => {
      jest.advanceTimersByTime(1000);
    });

    expect(screen.getByAltText("Slide 2")).toBeInTheDocument();
  });
});
