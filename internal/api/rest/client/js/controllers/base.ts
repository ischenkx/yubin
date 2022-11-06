import {Axios, AxiosResponse} from "axios";

class Base {
    constructor(protected client: Axios) {}

    protected async get(url: string, config?: any) {
        config = config || {}
        config.url = url
        config.method = 'GET'
        return this.request(config)
    }

    protected async post(url: string, data?: any, config?: any) {
        config = config || {}
        config.url = url
        config.data = data
        config.method = 'POST'
        return this.request(config)
    }

    protected async put(url: string, data?: any, config?: any) {
        config = config || {}
        config.url = url
        config.data = data
        config.method = 'PUT'
        return this.request(config)
    }

    protected async delete(url: string, data?: any, config?: any) {
        config = config || {}
        config.url = url
        config.data = data
        config.method = 'DELETE'
        return this.request(config)
    }

    protected async request(config: any) {
        let response = await this.client.request(config)
        if (response.status != 200) {
            throw new Error(`failed to get data: status code ${response.status}`)
        }
        return response.data
    }
}

export default Base