"use client"
import { useState, useEffect } from "react"
import { cn } from "@/lib/utils"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import axios from 'axios';
import { useToast } from "@/hooks/use-toast"
import { ToastAction } from "@/components/ui/toast"
import { useCookies } from "react-cookie";

interface LoginFormProps extends React.ComponentProps<"div"> {
    type: "signin" | "register" | "forgot";
}

// const axiosInstance = axios.create({
//   baseURL: "http://localhost:8080/api/v1",
// })

const axiosInstance = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // Ensure this matches your API base URL
  headers: {
    'Content-Type': 'application/json', // Set default headers if required
  },
});

export function LoginForm({
  type,
  className,
  ...props
}: LoginFormProps) {

  const [userToken, setUserToken] = useState(null);
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [retypedPassword, setRetypedPassword] = useState("");

  const [cookies, setCookie, removeCookie] = useCookies(["cookie-name"]);

  const { toast } = useToast()

  const handleUsernameChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setUsername(event.target.value);
  };
  const handlePasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setPassword(event.target.value);
  };

  const handleRetypedPasswordChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRetypedPassword(event.target.value);
  };

  const signinRequest = async (username: string, password: string) => {
    try{
        const url = '/auth/login';
        const { data } = await axiosInstance.post(
            url,
            {
              "username": username,
              "password": password
            }
        );
        setUserToken(data.Key.token);
        setCookie("jwt", data.Key.token);
        setCookie("refresh_token", data.Key.refresh_token);
        location.pathname = "/";
        // console.log(data.Key.token);
        // console.log(data.Key.refresh_token);

        toast({
          title: "Authentication Successful!",
          description: "You have successfully signed into your Meduza Team Server.",
          action: (
            <ToastAction altText="undo">Undo</ToastAction>
          ),
        })
    }
    catch(error){
        toast({
          title: "Invalid Credentials!",
          description: "You have either entered invalid credentials or the server is down. Please try again later.",
          action: (
            <ToastAction altText="undo">Undo</ToastAction>
          ),
        })
    }
  }

  const registerRequest = async (username: string, password: string) => {
    try{
        const url = '/auth/login';
        const { data } = await axiosInstance.post(
            url,
            {
              "username": username,
              "password": password
            }
        );
        setUserToken(data.Key.token);
        setCookie("jwt", data.Key.token);
        setCookie("refresh_token", data.Key.refresh_token);
        location.pathname = "/";
        // console.log(data.Key.token);
        // console.log(data.Key.refresh_token);

        toast({
          title: "Authentication Successful!",
          description: "You have successfully signed into your Meduza Team Server.",
          action: (
            <ToastAction altText="undo">Undo</ToastAction>
          ),
        })
    }
    catch(error){
        toast({
          title: "Error With Registration!",
          description: "You have either entered mismatching passwords or the server is down. Please try again later.",
          action: (
            <ToastAction altText="undo">Undo</ToastAction>
          ),
        })
    }
  }


  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card className="overflow-hidden">
        <CardContent className="grid p-0 md:grid-cols-2">
          <form className="p-6 md:p-8">
            <div className="flex flex-col gap-6">
              <div className="flex flex-col items-center text-center">
                <h1 className="text-2xl font-bold">{type === "signin" ? "Welcome back" : "Register account"}</h1>
                <p className="text-balance text-muted-foreground">
                  {/* Login to your Acme Inc account */}
                  {type === "signin" ? "Login to your Meduza Team Server" : "Register an Account for Meduza"}
                </p>
              </div>
              <div className="grid gap-2">
                <Label htmlFor="email">Username</Label>
                <Input
                  id="username"
                  type="username"
                  placeholder="winston.churchill"
                  value={username}
                  onChange={handleUsernameChange}
                  required
                />
              </div>
              <div className="grid gap-2">
                <div className="flex items-center">
                  <Label htmlFor="password">Password</Label>
                  <a
                    href="#"
                    className="ml-auto text-sm underline-offset-2 hover:underline"
                  >
                    
                    {type === "signin" ? "Forgot your password?" : ""}
                  </a>
                </div>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={handlePasswordChange}
                  required
                />
              </div>
              {type === "register" ?
                <div className="grid gap-2">
                  <div className="flex items-center">
                    <Label htmlFor="password">Re-Type Password</Label>
                    <a
                      href="#"
                      className="ml-auto text-sm underline-offset-2 hover:underline"
                    >
                    </a>
                  </div>
                  <Input
                    id="password"
                    type="password"
                    value={retypedPassword}
                    onChange={handleRetypedPasswordChange}
                    required
                  />
                </div>

                :

                <></>
              }
              <Button
                type="submit"
                className="w-full"
                onClick={(event) => {
                  event.preventDefault();
                  if (type === "signin") {
                    signinRequest(username, password);
                  }
                }}
              >
                {type === "signin" ? "Login" : "Register"}
              </Button>
              <div className="text-center text-sm">
                {type === "signin" ? "Don't have an account? " : "Have an account? "}
                <a href={type === "signin" ? "register" : "signin"} className="underline underline-offset-4">
                  {type === "signin" ? "Register" : "Log in"}
                </a>
              </div>
            </div>
          </form>
          <div className="relative hidden md:flex md:flex-row md:items-center md:justify-center p-0 m-0">
            <img
              src="/meduza.png"
              alt="Image"
              // className="absolute inset-0 h-[50%] w-[50%] object-cover "
              className="block overflow-visible h-[55%] w-[55%] object-cover"
            />
          </div>
        </CardContent>
      </Card>
      <div className="text-balance text-center text-xs text-muted-foreground [&_a]:underline [&_a]:underline-offset-4 hover:[&_a]:text-primary">
        By clicking continue, you agree to our <a href="#">Terms of Service</a>{" "}
        and <a href="#">Privacy Policy</a>.
      </div>
    </div>
  )
}
