package main

import (
        "fmt"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/elasticache"
        "time"
)

func main() {

        endpoints := []string{}
        client := elasticache.New(session.New(),&aws.Config{Region: aws.String("eu-west-1")})
        params := &elasticache.DescribeCacheClustersInput{
            CacheClusterId:    aws.String("eticket-cache"),
            MaxRecords:        aws.Int64(20),
            ShowCacheNodeInfo: aws.Bool(true),
        }

        go func(endpoints *[]string) {
                for {
                        *endpoints = nil
                        resp, err := client.DescribeCacheClusters(params)

                        if err == nil {
                                for _, cluster := range resp.CacheClusters {
                                        for size, node := range cluster.CacheNodes {
                                                if size == 0 {
                                                        if *cluster.CacheNodes[0].CacheNodeStatus == "available" {
                                                                *endpoints = append(*endpoints, fmt.Sprintf("%s:%d", *cluster.CacheNodes[0].Endpoint.Address, *cluster.CacheNodes[0].Endpoint.Port))
                                                        }
                                                } else {
                                                        if *node.CacheNodeStatus == "available" {
                                                                *endpoints = append(*endpoints, fmt.Sprintf("%s:%d", *node.Endpoint.Address, *node.Endpoint.Port))
                                                        }
                                                }
                                        }
                                }
                                time.Sleep(time.Second * 2)
                        }
                }
        }(&endpoints)

        for {
                time.Sleep(time.Second * 2)
                if endpoints != nil {
                        fmt.Println(endpoints);
                }
        }

}
