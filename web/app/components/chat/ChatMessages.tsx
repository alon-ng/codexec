import { useTranslation } from "react-i18next";
import type { ModelsChatMessage } from "~/api/generated/model";
import ChatMessage from "./ChatMessage";

export interface ChatMessagesProps {
    isLoading: boolean;
    isLoadingResponse: boolean;
    messages: ModelsChatMessage[];
}

export default function ChatMessages({ isLoading, messages, isLoadingResponse }: ChatMessagesProps) {
    const { t } = useTranslation();
    const greetingMessage: ModelsChatMessage = {
        uuid: "00000000-0000-0000-0000-000000000000",
        ts: new Date().toISOString(),
        role: "assistant",
        content: t("cody.greeting"),
        exercise_uuid: "00000000-0000-0000-0000-000000000000",
        user_uuid: "00000000-0000-0000-0000-000000000000",
    }

    return (
        <div className="flex flex-col flex-1 overflow-y-auto gap-2">
            {isLoading ? (
                <ChatMessage />
            ) : (
                <>
                    <ChatMessage message={greetingMessage} />
                    {messages.map((message) => (
                        <ChatMessage key={message.uuid} message={message} />
                    ))}
                    {isLoadingResponse && (
                        <ChatMessage isLoadingCodeResponse={isLoadingResponse} />
                    )}
                </>
            )}
        </div>
    );
}