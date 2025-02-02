"use client";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Trash2 } from "lucide-react";
import { useState } from "react";
import { AddUserDialog } from "./add-user-dialog";
import { mockUsers } from "./data";
import { DeleteUserDialog } from "./delete-user-dialog";
import { User } from "./types";

export function UserManagement() {
  const [users, setUsers] = useState<User[]>(mockUsers);
  const [userToDelete, setUserToDelete] = useState<User | null>(null);
  const [isAddUserOpen, setIsAddUserOpen] = useState(false);

  const isAdmin = true; // TODO

  const handleDeleteUser = (user: User) => {
    setUsers(users.filter((u) => u.id !== user.id));
    setUserToDelete(null);
  };

  const handleAddUser = (newUser: User) => {
    setUsers([...users, newUser]);
    setIsAddUserOpen(false);
  };

  return (
    <Card className="w-full max-w-4xl mx-auto border-none max-h-[50vh] overflow-y-scroll">
      <CardHeader className="flex flex-row items-center justify-between">
        <CardTitle className="text-2xl font-bold">Teams</CardTitle>
        {isAdmin && (
          <Button onClick={() => setIsAddUserOpen(true)}>Add User</Button>
        )}
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {users.map((user) => (
            <div
              key={user.id}
              className="flex items-center justify-between p-4 bg-secondary rounded-lg"
            >
              <div className="flex items-center space-x-4">
                <Avatar className="h-12 w-12">
                  <AvatarImage src={user.avatarUrl} alt={user.name} />
                  <AvatarFallback>
                    {user.name.slice(0, 2).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <h3 className="font-semibold">{user.name}</h3>
                  <p className="text-sm text-muted-foreground">{user.email}</p>
                </div>
              </div>
              {isAdmin && (
                <Button
                  variant={"destructive"}
                  size="icon"
                  onClick={() => setUserToDelete(user)}
                >
                  <Trash2 className="h-5 w-5 text-white" />
                </Button>
              )}
            </div>
          ))}
        </div>
      </CardContent>
      {userToDelete && (
        <DeleteUserDialog
          user={userToDelete}
          onConfirm={() => handleDeleteUser(userToDelete)}
          onCancel={() => setUserToDelete(null)}
        />
      )}
      <AddUserDialog
        isOpen={isAddUserOpen}
        onClose={() => setIsAddUserOpen(false)}
        onAddUser={handleAddUser}
      />
    </Card>
  );
}