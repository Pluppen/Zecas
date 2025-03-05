import {atom} from 'nanostores';

export type User = {
    name?: string,
    email?: string,
    image?: string,
    access_token?: string
}

export const user = atom<User>();