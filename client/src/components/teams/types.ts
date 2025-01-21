export enum UserRole {
    Admin = "Admin",
    User = "User",
  }
  
  export type User = {
    id: string;
    name: string;
    email: string;
    role: UserRole;
    avatarUrl: string;
  };