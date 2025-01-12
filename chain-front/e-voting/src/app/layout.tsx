import type {Metadata} from "next";
import {Geist, Geist_Mono} from "next/font/google";
import "./globals.css";
import {PrimeReactProvider} from "primereact/api";
import "primereact/resources/primereact.css";
import "primeflex/primeflex.css";
import "primeicons/primeicons.css";
import ReactQueryProvider from "@/utils/providers/react-query-provider";

const geistSans = Geist({
    variable: "--font-geist-sans",
    subsets: ["latin"],
});

const geistMono = Geist_Mono({
    variable: "--font-geist-mono",
    subsets: ["latin"],
});

export const metadata: Metadata = {
    title: "E-Voting",
};

export default function RootLayout({
                                       children,
                                   }: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en">
        <ReactQueryProvider>
            <PrimeReactProvider>
                <body className={`${geistSans.variable} ${geistMono.variable}`}>
                {children}
                </body>
            </PrimeReactProvider>
        </ReactQueryProvider>
        </html>
    );
}
