FROM node:20-alpine AS build

WORKDIR /app

COPY package.json package-lock.json* ./

RUN npm install

COPY . .
RUN npm run build

FROM node:20-alpine AS production

WORKDIR /app

COPY package.json package-lock.json* ./

RUN npm install --omit=dev

COPY --from=build /app/build ./build
COPY --from=build /app/.svelte-kit ./.svelte-kit

EXPOSE 3000

ENV NODE_ENV=production

USER 1000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

CMD ["node", "build"]