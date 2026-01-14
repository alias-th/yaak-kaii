package grpcserver

import (
	"context"
	"log"
	"time"
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
	payload := &types.CreateProductPayload{
		Name:        req.GetName(),
		UserId:      req.GetUserId(),
		Description: req.GetDescription(),
		CategoryId:  req.GetCategoryId(),
	}

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
			Slug:        p.Slug,
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

func (s *ProductGRPCServer) CreateProductVariant(ctx context.Context, req *pb.CreateProductVariantRequest) (*pb.CreateProductVariantResponse, error) {
	payload := &types.CreateProductVariantPayload{
		ProductID:  req.GetProductId(),
		Price:      req.GetPrice(),
		Stock:      req.GetStock(),
		Attributes: req.GetAttributes(),
	}

	variantID, err := s.productService.CreateProductVariant(ctx, payload)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProductVariantResponse{
		VariantId: variantID,
	}, nil

}

func (s *ProductGRPCServer) GetProductBySlug(ctx context.Context, req *pb.GetProductBySlugRequest) (*pb.GetProductBySlugResponse, error) {
	productDetail, err := s.productService.GetProductBySlug(ctx, req.GetShopId(), req.GetSlug())
	if err != nil {
		return nil, err
	}

	product := &pb.ProductDetails{
		Id:          productDetail.Product.ID.String(),
		ShopId:      productDetail.Product.ShopID.String(),
		Name:        productDetail.Product.Name,
		Description: productDetail.Product.Description,
		Status:      productDetail.Product.Status,
		CategoryId:  productDetail.Product.CategoryID.String(),
		CreatedAt:   productDetail.Product.CreatedAt.Format(time.RFC3339),
		Slug:        productDetail.Product.Slug,
	}
	axes := make([]*pb.AxisDTO, 0, len(productDetail.Axes))
	for _, a := range productDetail.Axes {
		axes = append(axes, &pb.AxisDTO{
			Key:      a.Key,
			Label:    a.Label,
			Type:     a.Type,
			Required: a.Required,
			Options:  a.Options,
			Order:    int32(a.Order),
		})
	}
	variants := make([]*pb.VariantDTO, 0, len(productDetail.Variants))
	for _, v := range productDetail.Variants {
		variants = append(variants, &pb.VariantDTO{
			Id:         v.ID.String(),
			Sku:        v.SKU,
			Price:      v.Price,
			Stock:      v.Stock,
			Attributes: v.Attributes,
			VariantKey: v.VariantKey,
		})
	}
	res := &pb.GetProductBySlugResponse{
		Product:  product,
		Images:   productDetail.Images,
		Axes:     axes,
		Variants: variants,
	}

	return res, nil
}
