import type { Metadata } from "next";
import { Eczar } from "next/font/google";
import "./globals.css";
import Navigation from "@/components/Navigation";

const eczar = Eczar({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  variable: "--font-eczar",
});

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
    <html lang="en" className={eczar.variable}>
      <body>
        <Navigation />
        {children}
      </body>
    </html>
  );
}
