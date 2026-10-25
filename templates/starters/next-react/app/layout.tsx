import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Project",
  description: "Replace this starter after the design contract is locked.",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
