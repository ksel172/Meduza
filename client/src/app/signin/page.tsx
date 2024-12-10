/**
 * v0 by Vercel.
 * @see https://v0.dev/t/Vef4HeOGRJI
 * Documentation: https://v0.dev/docs#integrating-generated-code-into-your-nextjs-app
 */
import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import Link from "next/link"
import { Button } from "@/components/ui/button"
import { Navbar } from "@/components/util/navbar/navbar";
import { Icons } from "@/components/icons"

export default function Component() {
  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen pb-0 mb-0 mt-10 gap-10 font-[family-name:var(--font-geist-sans)]">
        <Navbar />
        <div className="grid min-h-screen w-full grid-cols-1 lg:grid-cols-2">
            <div className="hidden bg-gray-100 lg:block dark:bg-gray-800">
                {/* <img
                src="/next.svg"
                alt="Login Illustration"
                width={1920}
                height={1080}
                className="h-full w-full object-cover"
                style={{ aspectRatio: "1920/1080", objectFit: "cover" }}
                /> */}
                <Icons.logo className="h-full w-full object-cover" style={{ aspectRatio: "1920/1080", objectFit: "cover" }}/>
            </div>
            <div className="flex items-center justify-center p-6 lg:p-10">
                <div className="mx-auto w-full max-w-[400px] space-y-6">
                <div className="space-y-2 text-center">
                    <h1 className="text-3xl font-bold">Welcome back!</h1>
                    <p className="text-gray-500 dark:text-gray-400">Enter your email and password to sign in.</p>
                </div>
                <form className="space-y-4">
                    <div className="space-y-2">
                    <Label htmlFor="email">Email</Label>
                    <Input id="email" type="email" placeholder="m@example.com" required />
                    </div>
                    <div className="space-y-2">
                    <div className="flex items-center justify-between">
                        <Label htmlFor="password">Password</Label>
                        <Link href="#" className="text-sm font-medium underline" prefetch={false}>
                        Forgot password?
                        </Link>
                    </div>
                    <Input id="password" type="password" required />
                    </div>
                    <Button type="submit" className="w-full">
                    Sign in
                    </Button>
                </form>
                </div>
            </div>
        </div>
    </div>
    
  )
}