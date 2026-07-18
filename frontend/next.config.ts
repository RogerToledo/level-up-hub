import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  rewrites: async () => [
    {
      source: "/api/:path*",
      destination: `${process.env.BACKEND_URL || "http://localhost:8081"}/v1/:path*`,
    },
  ],
};

export default nextConfig;
