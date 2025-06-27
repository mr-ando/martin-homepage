
import type { APIRoute } from "astro";
export const prerender = false; // keeps the endpoint dynamic, everything else is prerendered

export const GET: APIRoute = async ({ cookies, redirect }) => {
  cookies.delete("sb-access-token", { path: "/" });
  cookies.delete("sb-refresh-token", { path: "/" });
  return redirect("/signin");
};