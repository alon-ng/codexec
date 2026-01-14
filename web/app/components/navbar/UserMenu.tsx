import { LogOut } from "lucide-react";
import { useMeStore } from "~/stores/meStore";
import { usePostAuthLogout } from "~/api/generated/auth/auth";
import { useNavigate } from "react-router";
import { Avatar, AvatarFallback } from "~/components/ui/avatar";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";

export function UserMenu() {
    const navigate = useNavigate();
    const { firstName, lastName, email, clearMe } = useMeStore();
    const logoutMutation = usePostAuthLogout({
        mutation: {
            onSuccess: () => {
                clearMe();
                navigate("/login");
            },
            onError: () => {
                clearMe();
                navigate("/login");
            },
        },
    });

    const displayName = firstName;

    const initials =
        firstName && lastName
            ? `${firstName[0]}${lastName[0]}`.toUpperCase()
            : firstName
                ? firstName[0].toUpperCase()
                : lastName
                    ? lastName[0].toUpperCase()
                    : email
                        ? email[0].toUpperCase()
                        : "U";

    return (
        <DropdownMenu>
            <DropdownMenuTrigger asChild>
                <button className="group inline-flex items-center gap-2 cursor-pointer opacity-90 text-sm font-medium transition-opacity hover:bg-transparent hover:opacity-100 focus:bg-transparent focus:text-accent-foreground focus:outline-none">
                    <Avatar className="h-8 w-8">
                        <AvatarFallback className="bg-linear-65 from-purple-500 to-pink-500 text-primary-foreground text-sm font-medium">
                            {initials}
                        </AvatarFallback>
                    </Avatar>
                    <span className="hidden text-sm font-medium sm:inline-block">
                        {displayName}
                    </span>
                </button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
                <DropdownMenuLabel>
                    <div className="flex flex-col space-y-1">
                        <p className="text-sm font-medium leading-none">{displayName}</p>
                        {email && (
                            <p className="text-xs leading-none text-muted-foreground">{email}</p>
                        )}
                    </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuItem
                    onClick={() => logoutMutation.mutate()}
                    disabled={logoutMutation.isPending}
                    className="cursor-pointer text-destructive focus:text-destructive"
                >
                    <LogOut className="me-2 h-4 w-4" />
                    <span>{logoutMutation.isPending ? "Logging out..." : "Logout"}</span>
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    );
}
