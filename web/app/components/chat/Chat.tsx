import { Send } from 'lucide-react';
import { motion } from 'motion/react';
import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';
import { useGetMeExercisesExerciseUuidChat, usePostMeExercisesExerciseUuidChat } from '~/api/generated/me/me';
import type { ModelsChatMessage } from '~/api/generated/model';
import { Button } from '~/components/base/Button';
import ChatMessages from '~/components/chat/ChatMessages';
import { Textarea } from '~/components/ui/textarea';
import { blurInVariants } from '~/utils/animations';

export interface ChatProps {
    exerciseUuid: string;
    exerciseInstructions: string;
    exerciseCode: string;
}

export default function Chat({ exerciseUuid, exerciseInstructions, exerciseCode }: ChatProps) {
    const { t } = useTranslation();
    const [message, setMessage] = useState<string>("");
    const [messages, setMessages] = useState<ModelsChatMessage[]>([]);
    const { mutate, isPending } = usePostMeExercisesExerciseUuidChat({
        mutation: {
            onSuccess: (m) => {
                messages.push(m);
                setMessages([...messages]);
            },
            onError: (error) => {
                toast.error(error.error || t("chat.error"))
                const m = messages.pop();
                if (m) {
                    setMessage(m.content);
                }

                setMessages([...messages]);
            },
        },
    });

    const { data: messagesData, isLoading: isLoadingMessages } = useGetMeExercisesExerciseUuidChat(exerciseUuid);
    useEffect(() => {
        if (messagesData) {
            setMessages(messagesData);
        }
    }, [messagesData]);

    const handleSendMessage = () => {
        messages.push({
            uuid: "00000000-0000-0000-0000-000000000000",
            ts: new Date().toISOString(),
            role: "user",
            content: message,
            exercise_uuid: exerciseUuid,
            user_uuid: "00000000-0000-0000-0000-000000000000",
        })
        setMessages([...messages]);
        setMessage("");

        mutate({
            exerciseUuid, data: {
                content: message,
                code: exerciseCode,
                exercise_instructions: exerciseInstructions
            }
        });
    }

    return (
        <motion.div variants={blurInVariants()} initial="hidden" animate="visible" className="flex-1 flex flex-col gap-2 bg-muted rounded-lg overflow-hidden p-2">
            <ChatMessages messages={messages} isLoading={isLoadingMessages} isLoadingResponse={isPending} />
            <div className="relative rounded-lg shadow-sm bg-background p-2">
                <Textarea
                    className="bg-background resize-none max-h-48 min-h-0 h-8 focus-visible:ring-0 focus:border-0 border-0 shadow-none p-0 mb-10"
                    rows={1}
                    value={message}
                    onChange={(e) => setMessage(e.target.value)}
                    onKeyDown={(e) => {
                        if (e.key === 'Enter' && !e.shiftKey) {
                            e.preventDefault();
                            handleSendMessage();
                        }
                    }}
                />

                <Button
                    className="absolute right-2 bottom-2 rounded-full"
                    size="icon"
                    disabled={isLoadingMessages || !message || isPending}
                    onClick={handleSendMessage}
                >
                    <Send className="size-4" />
                </Button>
            </div>
        </motion.div>
    );
}