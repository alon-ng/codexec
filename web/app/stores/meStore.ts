import { create } from "zustand";
import { getMe } from "~/api/generated/me/me";

interface MeState {
    firstName: string | null;
    lastName: string | null;
    email: string | null;
    isLoggedIn: boolean;
    userUUID: string | null;
    isLoading: boolean;

    loadMe: () => Promise<void>;
    clearMe: () => void;
}

export const useMeStore = create<MeState>((set, get) => ({
    firstName: null,
    lastName: null,
    email: null,
    isLoggedIn: false,
    userUUID: null,
    isLoading: false,

    loadMe: async () => {
        // Prevent multiple simultaneous calls
        if (get().isLoading) {
            return;
        }

        set({ isLoading: true });

        try {
            const me = await getMe();
            set({
                firstName: me.first_name,
                lastName: me.last_name,
                email: me.email,
                userUUID: me.uuid,
                isLoggedIn: true,
                isLoading: false,
            });
        } catch (error: any) {
            get().clearMe();
        }
    },

    clearMe: () => {
        set({
            firstName: null,
            lastName: null,
            email: null,
            userUUID: null,
            isLoggedIn: false,
            isLoading: false,
        });
    },
}));
