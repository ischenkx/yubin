import Base from './base'
import {Axios} from "axios";
import Publication from "../models/publication";
import Report from "../models/report";
import Subscription from "../models/subscription";

export default class PubSub extends Base {
    constructor(client: Axios) {
        super(client);
    }

    async getTopics(offset = null, limit = null) {
        return await this
            .get('/pubsub/topics', {
                params: {offset, limit}
            })
            .then(response => response.data)
    }

    async getSubscribers(topicId: string, offset = null, limit = null) {
        return await this
            .get(`/pubsub/topics/${topicId}/subscribers`, {
                params: {offset, limit}
            })
            .then(response => response.data)
    }

    async deleteTopic(topicId: string) {
        await this.delete(`/pubsub/topics/${topicId}`)
    }

    async publish(publication: Publication) {
        return await this.post('/pubsub/publisher/publish', {
            source: publication.source,
            template: publication.template,
            users: publication.users,
            topics: publication.topics,
            meta: publication.meta,
            at: publication.at
        })
    }

    async getPublication(id: string) {
        return await this
            .get(`/pubsub/publisher/${id}`)
            .then(Publication.fromRaw)
    }

    async getReport(id: string) {
        return await this
            .get(`/pubsub/publisher/${id}/report`)
            .then(Report.fromRaw)
    }

    async getPersonalReport(id: string, userId: string) {
        return await this
            .get(`/pubsub/publisher/${id}/report/${userId}`)
            .then(Report.fromRaw)
    }

    async getReports(offset = null, limit = null) {
        return await this
            .get('/pubsub/publisher/reports', {
                params: {
                    offset, limit
                }
            }).
            then(data => (data || []).map(Report.fromRaw))
    }

    async getPublications(offset = null, limit = null) {
        return await this
            .get('/pubsub/publisher', {
                params: {
                    offset, limit
                }
            }).
            then(data => (data || []).map(Publication.fromRaw))
    }

    async getUserSubscriptions(userId: string) {
        return await this
            .get(`/pubsub/subscriptions/${userId}`)
            .then(data => (data || []).map(Subscription.fromRaw))
    }

    async getSubscription(userId: string, topic: string) {
        return await this
            .get(`/pubsub/subscriptions/${userId}/${topic}`)
            .then(Subscription.fromRaw)
    }

    async subscribe(userId: string, topic: string) {
        return await this
            .post(`/pubsub/subscriptions/${userId}/${topic}`)
    }

    async unsubscribe(userId: string, topic: string) {
        return await this
            .delete(`/pubsub/subscriptions/${userId}/${topic}`)
    }

    async unsubscribeAll(userId: string) {
        return await this
            .delete(`/pubsub/subscriptions/${userId}`)
    }
}