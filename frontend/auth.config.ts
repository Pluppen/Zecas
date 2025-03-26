import GitHub from '@auth/core/providers/github';
import Sendgrid from '@auth/core/providers/sendgrid';
import {loadEnv } from "vite";
import { defineConfig } from 'auth-astro';
import PostgresAdapter from "@auth/pg-adapter"
import pg from 'pg';

const { DATABASE_HOST, DATABASE_NAME, DATABASE_PASSWORD, DATABASE_USER } = loadEnv(process.env.NODE_ENV ?? "", process.cwd(), "");

const pool = new pg.Pool({
  host: DATABASE_HOST,
  user: DATABASE_USER,
  password: DATABASE_PASSWORD,
  database: DATABASE_NAME,
  max: 20,
  idleTimeoutMillis: 30000,
  connectionTimeoutMillis: 2000,
})

export default defineConfig({
    providers: [
        Sendgrid({
          apiKey: import.meta.env.SENDGRID_API_KEY,
          from: import.meta.env.EMAIL_AUTH_FROM
        }),
        GitHub({
          clientId: import.meta.env.GITHUB_CLIENT_ID,
          clientSecret: import.meta.env.GITHUB_CLIENT_SECRET,
        }),
    ],
    adapter: PostgresAdapter(pool),
})
