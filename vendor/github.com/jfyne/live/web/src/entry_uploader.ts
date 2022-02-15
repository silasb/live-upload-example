import { Socket } from "./socket";
import { LiveEvent } from "./event";

export interface UploadEntry {
    file: any,
    field: any,
    progress: any,
    done: any,
}

class Channel {
    private socket: Socket;
    private name: string;
    private entry: UploadEntry;

    constructor(socket: Socket, name: string, entry: UploadEntry) {
        this.socket = socket
        this.name = name
        this.entry = entry
    }

    public serialize() {
        this.entry.file.toJSON = () => {
            return {
                'lastModified': this.entry.file.lastModified,
                'lastModifiedDate': this.entry.file.lastModifiedDate,
                'name': this.entry.file.name,
                'size': this.entry.file.size,
                'type': this.entry.file.type
            }
        }
        return this.entry.file
    }

    public push(chunk: ArrayBuffer) {
        const base64String = btoa(String.fromCharCode(...new Uint8Array(chunk)));
        const data = {
            file: this.serialize(),
            chunk: base64String,
            field: this.entry.field,
        }

        const e = new LiveEvent(this.name, data, LiveEvent.GetID())
        return Socket.push(e)
    }
}

export class EntryUploader {
    private liveSocket: Socket;
    private entry: UploadEntry;
    private offset: number = 0;
    private chunkSize: number;
    private chunkTimer: number | null;
    private uploadChannel: Channel;

    constructor(entry, chunkSize, liveSocket) {
        this.liveSocket = liveSocket
        this.entry = entry
        this.offset = 0
        this.chunkSize = chunkSize
        this.chunkTimer = null
        this.uploadChannel = new Channel(liveSocket, 'allow_upload', entry)
        //this.uploadChannel = liveSocket.channel(`lvu:${entry.ref}`, {token: entry.metadata()})
    }

    error(reason) {
        clearTimeout(this.chunkTimer!)
        //this.uploadChannel.leave()
        //this.entry.error(reason)
    }

    upload() {
        setTimeout(() => {
            this.readNextChunk()
        }, 200)

        //this.uploadChannel.onError(reason => this.error(reason))
        //this.uploadChannel.join()
        //.receive("ok", _data => this.readNextChunk())
        //.receive("error", reason => this.error(reason))
    }

    isDone() { return this.offset >= this.entry.file.size }

    readNextChunk() {
        let reader = new window.FileReader()
        let blob = this.entry.file.slice(this.offset, this.chunkSize + this.offset)
        reader.onload = (e) => {
            if (e?.target?.error === null) {
                const chunk = e?.target?.result as ArrayBuffer
                this.offset += chunk?.byteLength
                this.pushChunk(chunk)
            } else {
                return console.log("Read error: " + e?.target?.error)
            }
        }
        reader.readAsArrayBuffer(blob)
    }

    pushChunk(chunk: ArrayBuffer) {
        this.uploadChannel.push(chunk)
            .receive("ok", () => {
                this.entry.progress((this.offset / this.entry.file.size) * 100)
                if (!this.isDone()) {
                    this.chunkTimer = window.setTimeout(() => this.readNextChunk(), 0)
                } else {
                    this.entry.done()
                }
            })
        //if(!this.uploadChannel.isJoined()){ return }
        //this.uploadChannel.push("chunk", chunk)
        //.receive("ok", () => {
        //if(!this.isDone()){
        //this.chunkTimer = setTimeout(() => this.readNextChunk(), this.liveSocket.getLatencySim() || 0)
        //}
        //})
    }
}
