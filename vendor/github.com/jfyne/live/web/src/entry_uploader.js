import { LiveEvent } from "./event";

class Channel {
	constructor(socket, name, entry) {
		this.socket = socket
		this.name = name
		this.entry = entry
	}


	serialize() {
		this.entry.file.toJSON = () => {
			return {
				'lastModified'     : this.entry.file.lastModified,
				'lastModifiedDate' : this.entry.file.lastModifiedDate,
				'name'             : this.entry.file.name,
				'size'             : this.entry.file.size,
				'type'             : this.entry.file.type 
			}
		}
		return this.entry.file
	}

	push(chunk) {
		const base64String = btoa(String.fromCharCode(...new Uint8Array(chunk)));
		const data = {
			file: this.serialize(),
			chunk: base64String,
			field: this.entry.field,
		}

		const e = new LiveEvent(this.name, data, LiveEvent.GetID())
		return this.socket.push(e)
	}
}

export default class EntryUploader {
  constructor(entry, chunkSize, liveSocket){
    this.liveSocket = liveSocket
    this.entry = entry
    this.offset = 0
    this.chunkSize = chunkSize
    this.chunkTimer = null
		this.uploadChannel = new Channel(liveSocket, 'allow_upload', entry)
		// this.onComplete = this.done
		//this.uploadChannel = liveSocket.channel(`lvu:${entry.ref}`, {token: entry.metadata()})
  }

  error(reason){
    clearTimeout(this.chunkTimer)
    //this.uploadChannel.leave()
    //this.entry.error(reason)
  }

  upload(){
		//this.liveSocket.

		setTimeout(() => {
			this.readNextChunk()
		}, 200)

    //this.uploadChannel.onError(reason => this.error(reason))
    //this.uploadChannel.join()
      //.receive("ok", _data => this.readNextChunk())
      //.receive("error", reason => this.error(reason))
  }

  isDone(){ return this.offset >= this.entry.file.size }

  readNextChunk(){
    let reader = new window.FileReader()
    let blob = this.entry.file.slice(this.offset, this.chunkSize + this.offset)
    reader.onload = (e) => {
      if(e.target.error === null){
        this.offset += e.target.result.byteLength
        this.pushChunk(e.target.result)
      } else {
        return console.log("Read error: " + e.target.error)
      }
    }
    reader.readAsArrayBuffer(blob)
  }

  pushChunk(chunk){
		this.uploadChannel.push(chunk)
			.receive("ok", () => {
				this.entry.progress((this.offset / this.entry.file.size) * 100)
				if(!this.isDone()){
					this.chunkTimer = setTimeout(() => this.readNextChunk(), 0)
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
