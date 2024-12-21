import { User, UserRole } from "./types";

export const mockUsers: User[] = [
  {
    id: "1",
    name: "John Doe",
    email: "john@example.com",
    role: UserRole.Admin,
    avatarUrl: "https://api.dicebear.com/6.x/avataaars/svg?seed=John",
  },
  {
    id: "2",
    name: "Jane Smith",
    email: "jane@example.com",
    role: UserRole.User,
    avatarUrl: "https://api.dicebear.com/6.x/avataaars/svg?seed=Jane",
  },
  {
    id: "3",
    name: "Bob Johnson",
    email: "bob@example.com",
    role: UserRole.User,
    avatarUrl: "https://api.dicebear.com/6.x/avataaars/svg?seed=Bob",
  },
];