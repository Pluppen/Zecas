import { defineMiddleware } from "astro/middleware";

import { getSession } from "auth-astro/server";

export const onRequest = defineMiddleware(async (context, next) => {
  // Add a string value to the locals object
  const session = await getSession(context.request)
  const isAuthed = session && "user" in session

  if(context.originPathname.startsWith("/api/auth")) {
    return next()
  }

  if(["/login"].includes(context.originPathname) && !isAuthed) {
    return next();
  }

  if(isAuthed && context.originPathname == "/login") {
    return context.redirect("/")
  }

  if (isAuthed) {
    return next();
  }

  return context.redirect("/login")

}); 