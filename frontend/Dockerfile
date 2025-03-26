FROM node:lts-alpine AS base
WORKDIR /app

COPY . .

RUN npm install --force
RUN npm run build -- --mode=docker

FROM node:lts-alpine AS runtime
WORKDIR /app

COPY --from=base /app/node_modules/ ./node_modules
COPY --from=base /app/dist/ ./dist

RUN ls -R ./dist/

ENV HOST=0.0.0.0
ENV PORT=4321
EXPOSE 4321
CMD node ./dist/server/entry.mjs
