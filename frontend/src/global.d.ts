/// <reference types="svelte" />

declare type ResponseData = {
    success: boolean
    result?: string
    error?: string
}

declare type Route = {
    [key: string]: boolean
}

declare type User = {
    id?: number
    name?: string
    email?: string
}

declare type Playload = {
    user: User
    iat: number
}