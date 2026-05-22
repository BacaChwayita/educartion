"use client";
import React from "react";
import Slider from "../components/Slider";
import Categories from "../components/Categories";
import FeaturedProducts from "../components/FeaturedProducts";

const slides = [
	{ id: 1, image: "/images/slide1.jpg", title: "Spring Collection", subtitle: "New arrivals" },
	{ id: 2, image: "/images/slide2.jpg", title: "Top Deals", subtitle: "Save up to 50%" },
];

const categories = [
	{ id: "cat-1", name: "Clothing", image: "/images/cat-clothing.jpg" },
	{ id: "cat-2", name: "Electronics", image: "/images/cat-electronics.jpg" },
	{ id: "cat-3", name: "Home", image: "/images/cat-home.jpg" },
	{ id: "cat-4", name: "Toys", image: "/images/cat-toys.jpg" },
	{ id: "cat-5", name: "Sports", image: "/images/cat-sports.jpg" },
	{ id: "cat-6", name: "Books", image: "/images/cat-books.jpg" },
	{ id: "cat-7", name: "Beauty", image: "/images/cat-beauty.jpg" },
	{ id: "cat-8", name: "Food", image: "/images/cat-food.jpg" },
];

const featured = [
	{ id: "p-1", title: "Comfort Tee", price: 19.99, image: "/images/product1.jpg", description: "Soft cotton tee" },
	{ id: "p-2", title: "Noise-cancelling Headphones", price: 129.99, image: "/images/product2.jpg", description: "Great sound" },
	{ id: "p-3", title: "Ceramic Mug", price: 12.0, image: "/images/product3.jpg", description: "Handmade" },
	{ id: "p-4", title: "Wireless Mouse", price: 34.99, image: "/images/product4.jpg", description: "Ergonomic design" },
	{ id: "p-5", title: "Running Shoes", price: 89.99, image: "/images/product5.jpg", description: "Professional grade" },
	{ id: "p-6", title: "Desk Lamp", price: 45.0, image: "/images/product6.jpg", description: "LED technology" },
	{ id: "p-7", title: "Water Bottle", price: 24.99, image: "/images/product7.jpg", description: "Stainless steel" },
	{ id: "p-8", title: "Yoga Mat", price: 29.99, image: "/images/product8.jpg", description: "Non-slip surface" },
	{ id: "p-9", title: "Portable Charger", price: 49.99, image: "/images/product9.jpg", description: "Fast charging" },
	{ id: "p-10", title: "Phone Stand", price: 15.99, image: "/images/product10.jpg", description: "Adjustable angle" },
];

export default function Home() {
	const addToCart = (product: any) => {
		// Simple placeholder: developer can replace with real cart integration
		const current = typeof window !== "undefined" ? JSON.parse(window.localStorage.getItem("cart") || "[]") : [];
		current.push({ ...product, qty: 1 });
		if (typeof window !== "undefined") window.localStorage.setItem("cart", JSON.stringify(current));
		console.log("Added to cart:", product);
		alert(`${product.title} added to cart`);
	};

	const onCategorySelect = (id: string | number) => {
		console.log("Category selected:", id);
	};

	return (
		<main style={{ width: "100%", display: "grid", gap: 28 }}>
			<div style={{ width: "100vw", position: "relative", left: "50%", right: "50%", marginLeft: "-50vw", marginRight: "-50vw" }}>
				<Slider slides={slides} />
			</div>

		<div style={{ padding: "0 20px", maxWidth: 1200, margin: "0 auto", display: "grid", gap: 28 }}>
				<section>
					<h2 style={{ marginBottom: 12 }}>Categories</h2>
					<Categories categories={categories} onSelect={onCategorySelect} />
				</section>

				<section>
					<h2 style={{ marginBottom: 12 }}>Featured Products</h2>
					<FeaturedProducts products={featured} onAddToCart={addToCart} />
				</section>
			</div>
		</main>
	);
}
