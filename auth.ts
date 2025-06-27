import { betterAuth } from "better-auth";
import { Pool } from "pg";

export const auth = betterAuth({
    //...
    database: new Pool({
      connectionString: process.env.DATABASE_URL || "postgres://user:password@localhost:5432/mydb"
    }),
    socialProviders: {
        github: {
          clientId: process.env.GITHUB_CLIENT_ID || "",
          clientSecret: process.env.GITHUB_CLIENT_SECRET || "",
        }
    }
})