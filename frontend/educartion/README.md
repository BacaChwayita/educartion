# Educartion Frontend

This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Overview

The frontend module provides the user interface for the Educartion application. It is built with Next.js and communicates with the backend API to fetch and manage educational content, products, orders, and user data.

## Tech Stack

- **Framework**: [Next.js](https://nextjs.org) (React-based)
- **Package Manager**: npm, yarn, pnpm, or bun
- **Styling**: CSS Modules and Tailwind CSS
- **Font Optimization**: [Geist](https://vercel.com/font) font family

## Frontend Structure

```text
frontend/educartion/
├── src/
│   ├── app/                 # Next.js App Router pages and layouts
│   │   ├── api/             # API routes and handlers
│   │   ├── admin/           # Admin dashboard pages
│   │   ├── catalog/         # Product catalog pages
│   │   ├── cart/            # Shopping cart pages
│   │   ├── products/        # Product detail pages
│   │   ├── payment/         # Payment and checkout pages
│   │   ├── login/           # Authentication pages
│   │   ├── register/        # User registration pages
│   │   ├── layout.tsx       # Root layout
│   │   └── page.tsx         # Home page
│   ├── components/          # Reusable React components
│   │   ├── Categories.tsx
│   │   ├── FeaturedProducts.tsx
│   │   ├── Slider.tsx
│   │   ├── admin-dashboard.tsx
│   │   ├── auth-form.tsx
│   │   └── ...
│   ├── controllers/         # Business logic controllers
│   │   ├── adminController.ts
│   │   ├── authController.ts
│   │   ├── cartController.ts
│   │   ├── paymentController.ts
│   │   └── productController.ts
│   ├── services/            # API communication services
│   │   ├── authService.ts
│   │   ├── cartService.ts
│   │   ├── paymentService.ts
│   │   └── productService.ts
│   ├── lib/                 # Shared utilities and contracts
│   │   ├── auth-contract.ts
│   │   ├── auth-errors.ts
│   │   ├── backend-auth.ts
│   │   ├── cart-contract.ts
│   │   ├── payment-contract.ts
│   │   └── product-contract.ts
│   └── globals.css          # Global styles
├── public/                  # Static assets
├── package.json             # Dependencies and scripts
├── next.config.ts           # Next.js configuration
├── tsconfig.json            # TypeScript configuration
├── eslint.config.mjs        # ESLint configuration
└── postcss.config.mjs       # PostCSS configuration
```

## Getting Started

### Prerequisites

- Node.js (v18 or higher)
- npm, yarn, pnpm, or bun

### Setup

1. Navigate to the frontend directory:

```bash
cd frontend/educartion
```

2. Install dependencies:

```bash
npm install
# or
yarn install
# or
pnpm install
# or
bun install
```

### Run Development Server

Start the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

The page auto-updates as you edit the file. Start by modifying `src/app/page.tsx` to see changes in real-time.

### Build for Production

Create an optimized production build:

```bash
npm run build
# or
yarn build
# or
pnpm build
# or
bun build
```

### Run Production Build

After building, you can start the production server:

```bash
npm run start
# or
yarn start
# or
pnpm start
# or
bun start
```

## Key Features

### Pages

- **Home** (`/`) - Landing page with featured products and categories
- **Catalog** (`/catalog`) - Browse all products
- **Products** (`/products`) - Product detail pages
- **Cart** (`/cart`) - Shopping cart management
- **Payment** (`/payment`) - Checkout and payment processing
- **Login** (`/login`) - User authentication
- **Register** (`/register`) - User account creation
- **Admin** (`/admin`) - Admin dashboard for managing content

### API Routes

- `POST /api/auth/register` - User registration
- `POST /api/auth/login` - User login

### Controllers & Services

The frontend uses a controller-service architecture:

- **Controllers** - Handle business logic and coordinate between components and services
- **Services** - Manage API communication with the backend
- **Contracts** - Type definitions for API request/response contracts

## Environment Variables

Create a `.env.local` file in the `frontend/educartion` directory with:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080/api
```

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - Learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - An interactive Next.js tutorial.
- [Next.js GitHub Repository](https://github.com/vercel/next.js) - Your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out the [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.

## Related Documentation

- Backend API documentation: `../../backend/README.md`
- Database model (ERD source): `../../data/ERD.mmd`
- Database setup scripts: `../../data/scripts/`

## Development Guidelines

- Follow TypeScript best practices
- Use ESLint for code quality
- Create feature branches for new development
- Ensure all changes are tested before submitting pull requests
