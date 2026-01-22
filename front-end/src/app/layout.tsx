import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Pizza Vibe",
  description: "Order delicious pizzas",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
