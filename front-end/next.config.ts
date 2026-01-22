import type { NextConfig } from "next";

const storeUrl = process.env.STORE_SERVICE_URL || 'http://localhost:8080';

const nextConfig: NextConfig = {
  output: 'standalone',
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: `${storeUrl}/:path*`,
      },
    ];
  },
};

export default nextConfig;
