FROM node:22

RUN corepack enable pnpm

WORKDIR /app

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./

RUN pnpm install --frozen-lockfile

COPY . .

EXPOSE 3000

# `pnpm run dev` is dev-only server. For prod we need build + `next start`
CMD ["pnpm", "run", "dev"]
