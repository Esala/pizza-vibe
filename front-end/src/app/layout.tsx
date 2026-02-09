import type { Metadata } from "next";
import { Geist } from "next/font/google";
import { Knewave } from "next/font/google";
import "./tokens.css";
import "./globals.css";
import Navigation from "@/components/Navigation";

const geist = Geist({
  subsets: ["latin"],
  variable: "--font-geist",
});

const knewave = Knewave({
  weight: "400",
  subsets: ["latin"],
  variable: "--font-knewave",
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
    <html lang="en" className={`${geist.variable} ${knewave.variable}`}>
      <body>
        <Navigation />
        {children}
      </body>
    </html>
  );
}
