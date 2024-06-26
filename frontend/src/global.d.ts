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
    lastEditTime?: number
    registerTime?: number
    deletedTime?: number
}