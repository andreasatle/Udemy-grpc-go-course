package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/andreasatle/Udemy/grpc-go-course/blog/pb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := pb.NewBlogServiceClient(cc)

	blog := &pb.Blog{
		AuthorId: "Andreas",
		Title:    "My first blog",
		Content:  "Content of the first blog",
	}
	blogId := createBlog(c, blog)
	fmt.Println(blogId)
	blog1 := readBlog(c, "5f46af2cdcde90d5710ed301")
	fmt.Println(blog1)
	blog2 := readBlog(c, blogId)
	fmt.Println(blog2)
	newBlog := &pb.Blog{
		Id:       blogId,
		AuthorId: "New Author",
		Title:    "Wind of Change",
		Content:  "Change, change",
	}
	updateBlog(c, newBlog)
	//deleteBlog(c, blogId)
	listBlog(c)
}

func createBlog(c pb.BlogServiceClient, blog *pb.Blog) string {
	fmt.Printf("Creating a blog: %v\n", blog)

	res, err := c.CreateBlog(context.Background(), &pb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}

	fmt.Printf("Blog has been created: %v\n", res)

	return res.Blog.GetId()
}

func readBlog(c pb.BlogServiceClient, blogId string) *pb.Blog {
	fmt.Printf("Reading a blog: %v\n", blogId)

	fmt.Printf("BlogId: %v\n", blogId)
	req := &pb.ReadBlogRequest{BlogId: blogId}
	res, err := c.ReadBlog(context.Background(), req)
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
		return nil
	}

	fmt.Printf("Blog has been read: %v\n", res)
	return res.Blog
}

func updateBlog(c pb.BlogServiceClient, blog *pb.Blog) {
	fmt.Printf("Updating a blog: %v\n", blog)
	res, err := c.UpdateBlog(context.Background(), &pb.UpdateBlogRequest{Blog: blog})
	if err != nil {
		fmt.Printf("Error happened while updating: %v\n", err)
		return
	}
	fmt.Printf("Blog was updated: %v\n", res)
}

func deleteBlog(c pb.BlogServiceClient, blogId string) {
	fmt.Printf("Deleting a blog: %v\n", blogId)

	req := &pb.DeleteBlogRequest{BlogId: blogId}
	res, err := c.DeleteBlog(context.Background(), req)
	if err != nil {
		fmt.Printf("Error happened while deleting: %v\n", err)
		return
	}

	fmt.Printf("Blog has been deleted: %v\n", res)
}

func listBlog(c pb.BlogServiceClient) {
	stream, err := c.ListBlog(context.Background(), &pb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Error calling ListBlog RPC: %v\n", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something went wrong: %v\n", err)
		}
		fmt.Println(res.GetBlog())
	}

}
