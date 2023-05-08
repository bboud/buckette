import { atom } from 'jotai'
import { DropzoneInputProps } from 'react-dropzone'

export const fileStatusAtom = atom<FileStatus[]>([])

export const isFocusedAtom = atom(false)
export const isDragAcceptAtom = atom(false)
export const isDragRejectAtom = atom(false)

export const inputPropsAtom = atom<() => DropzoneInputProps>(() => ({}))

export const filesAtom = atom<File[]>([])

export const isUploadingAtom = atom(false)
