import {atom} from 'nanostores';
import {persistentAtom} from '@nanostores/persistent';

export type Project = {
    id: string,
    name: string,
    description: string,
    targets: {
        ip_ranges: string,
        cidr_ranges: string,
        domains: string,
    }[]
}

export type ProjectData = {
    projects: Project[]
}

export const projects = atom<ProjectData>({
    projects: []
});

export const activeProjectStore = atom<Project | undefined>();

export const activeProjectIdStore = persistentAtom('activeProjectId')