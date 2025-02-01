"use client"

import { LoginForm } from "@/components/login-form"
import { useEffect } from "react"
import { useCookies } from "react-cookie";

export default function LoginPage() {

  const [cookies, setCookie, removeCookie] = useCookies(["cookie-name"]);

  useEffect(() => {
    if(cookies.jwt){
      window.location.pathname = "/";
    }
  })

  return (
    <div className="flex min-h-svh flex-col w-[100vw] items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm md:max-w-3xl">
        <LoginForm type="signin" />
      </div>
    </div>
  )
}
