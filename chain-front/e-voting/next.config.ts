import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async redirects() {
      return [
          {
              source: '/',
              destination: '/election',
              permanent: false,
          }
      ]
  }
};

export default nextConfig;
