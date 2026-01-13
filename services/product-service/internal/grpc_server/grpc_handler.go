package grpcserver

import (
	"context"
	"log"
	"yaak-kaii/services/product-service/internal/services"
	"yaak-kaii/services/product-service/pkg/types"
	pb "yaak-kaii/shared/proto/product"
)

type ProductGRPCServer struct {
	pb.UnimplementedProductServiceServer
	productService  *services.ProductService
	categoryService *services.CategoryService
}

func NewProductGRPCServer(product *services.ProductService, category *services.CategoryService) *ProductGRPCServer {
	return &ProductGRPCServer{
		productService:  product,
		categoryService: category,
	}
}

func (s *ProductGRPCServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	payload := &services.CreateProductPayload{
		Name:        req.GetName(),
		UserId:      req.GetUserId(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		CategoryId:  req.GetCategoryId(),
		Stock:       req.GetStock(),
		Attributes:  req.GetAttributes(),
	}

	category, err := s.categoryService.GetCategoryByID(ctx, req.GetCategoryId())
	if err != nil {
		return nil, err
	}
	categoryName := category.Name
	payload.CategoryName = categoryName

	productId, err := s.productService.CreateProduct(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProductResponse{
		ProductId: productId,
	}, nil
}

func (s *ProductGRPCServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	payload := types.ListProductsReq{
		Q:           req.GetQ(),
		CategoryIds: req.GetCategoryIds(),
		MinPrice:    req.GetMinPrice(),
		MaxPrice:    req.GetMaxPrice(),
		Limit:       req.GetLimit(),
		Page:        req.GetPage(),
		Sort:        req.GetSort(),
	}

	products, err := s.productService.ListProducts(ctx, payload)
	if err != nil {
		return nil, err
	}

	productRes := []*pb.Product{}
	for _, p := range products.Products {
		variants := []*pb.ProductVariant{}
		for _, v := range p.Variants {
			att, err := v.Attributes.Value()
			if err != nil {
				log.Println(err)
				return nil, err
			}

			variants = append(variants, &pb.ProductVariant{
				Sku:        v.Sku,
				Price:      v.Price,
				Stock:      int32(v.Stock),
				Attributes: att,
			})
		}
		images := []string{}
		for _, img := range p.Images {
			images = append(images, img.ImageURL)
		}

		productRes = append(productRes, &pb.Product{
			Id:          p.ID.String(),
			Name:        p.Name,
			ShopId:      p.ShopID.String(),
			Description: p.Description,
			Status:      p.Status,
			Variants:    variants,
			Images:      images,
			Category: &pb.ProductCategory{
				CategoryId: p.Category.ID.String(),
				Name:       p.Category.Name,
			},
		})
	}

	nextPage := int32(0)
	prevPage := int32(0)
	if products.Pagination.NextPage != nil {
		nextPage = int32(*products.Pagination.NextPage)
	}
	if products.Pagination.PrevPage != nil {
		prevPage = int32(*products.Pagination.PrevPage)
	}

	return &pb.ListProductsResponse{
		Products: productRes,
		Pagination: &pb.Pagination{
			Limit:    int32(products.Pagination.PageSize),
			Page:     int32(products.Pagination.Page),
			Total:    int32(products.Pagination.Total),
			LastPage: int32(products.Pagination.LastPage),
			NextPage: nextPage,
			PrevPage: prevPage,
			HasNext:  products.Pagination.HasNext,
			HasPrev:  products.Pagination.HasPrev,
		},
	}, nil
}

func (s *ProductGRPCServer) GetProductByID(ctx context.Context, req *pb.GetProductByIDRequest) (*pb.GetProductByIDResponse, error) {
	product, err := s.productService.GetProductByID(ctx, req.GetProductId())
	if err != nil {
		return nil, err
	}
	return &pb.GetProductByIDResponse{
		Product: &pb.Product{
			Id:          product.ID.String(),
			Name:        product.Name,
			ShopId:      product.ShopID.String(),
			Description: product.Description,
			Status:      product.Status,
			Category: &pb.ProductCategory{
				CategoryId: product.Category.ID.String(),
				Name:       product.Category.Name,
			},
		},
	}, nil
}

func (s *ProductGRPCServer) UploadProductImages(ctx context.Context, req *pb.UploadProductImagesRequest) (*pb.UploadProductImagesResponse, error) {
	err := s.productService.UploadProductImages(ctx, req.GetProductId(), req.GetImageUrls())
	if err != nil {
		return nil, err
	}

	return &pb.UploadProductImagesResponse{
		Success: true,
	}, nil
}
