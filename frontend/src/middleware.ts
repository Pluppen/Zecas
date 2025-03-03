import { defineMiddleware } from "astro/middleware";

import { getSession } from "auth-astro/server";

export const onRequest = defineMiddleware(async (context, next) => {
  // Add a string value to the locals object
  const session = await getSession(context.request)

  console.log(session)

  if(["/login", "/api/auth/csrf", "/api/auth/signin/github", "/api/auth/callback/github"].includes(context.originPathname)) {
    return next();
  }

  if (session && "user" in session) {
    return next();
  }

  return context.redirect("/login")

}); 