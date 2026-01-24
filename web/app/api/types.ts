export interface UserExerciseQuizData {
    [key: string]: string;
}

export interface UserExerciseCodeData {
    name: string;
    content?: string;
    children?: UserExerciseCodeData[];
}

export interface CheckerResult {
    type: string;
    success: boolean;
    message: string;
}

export interface ExecuteResponse {
    job_id: string;
    stdout: string;
    stderr: string;
    exit_code: number;
    time: number;
    memory: number;
    cpu: number;
    checker_results: CheckerResult[];
    passed: boolean;
    next_lesson_uuid?: string;
    next_exercise_uuid?: string;
}