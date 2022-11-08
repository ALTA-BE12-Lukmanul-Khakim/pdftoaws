package modivy

import (
	"bytes"
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
)

//const path = "gopdf/go-wkhtmltopdf-master"

// RequestPDF struct
type RequestPDF struct {
	body []byte
}

type Product struct {
	Item     string
	Price    float64
	Qty      float64
	Subtotal float64
}

type Data struct {
	Product       []Product
	Image         string
	InvoiceNo     int
	Date          string
	DueDate       string
	DeliveryDate  string
	PaymentMethod string
}

func CreateInv() (string, error) {

	//wkhtmltopdf.SetPath(path)

	var (
		products []Product
	)
	dt := time.Now()
	data := Data{
		Image:         "bengcall/merdeka.png",
		InvoiceNo:     1,
		Date:          dt.Format("01-02-2006"),
		DueDate:       dt.Format("01-02-2006"),
		DeliveryDate:  dt.Format("01-02-2006"),
		PaymentMethod: "Paypal",
	}
	products = append(products, Product{Item: "item", Price: 10000, Qty: 1, Subtotal: 10000})
	products = append(products, Product{Item: "item2", Price: 20000, Qty: 2, Subtotal: 40000})
	products = append(products, Product{Item: "item3", Price: 30000, Qty: 3, Subtotal: 90000})

	data.Product = products
	r := &RequestPDF{}
	err := r.ParseTemplate("invoice.html", data)

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(r.body)))
	pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile("./test.pdf")
	if err != nil {
		log.Fatal(err)
	}

	// Upload
	res, err := UploudFile(pdfg.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	return res, nil
}

// func AddFileToS3(s *session.Session, buffer []byte) (string, error) {
// 	fileName := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
// 	// Config settings: this is where you choose the bucket, filename, content-type etc.
// 	// of the file you're uploading.
// 	_ , err := s3.New(s).PutObject(&s3.PutObjectInput{
// 		Bucket:               aws.String("ap-southeast-1"),
// 		Key:                  aws.String("profile/" + fileName),
// 		ACL:                  aws.String("public-read"),
// 		Body:                 bytes.NewReader(buffer),
// 		ContentLength:        aws.Int64(int64(len(buffer))),
// 		ContentType:          aws.String(http.DetectContentType(buffer)),
// 		ContentDisposition:   aws.String("attachment"),
// 		ServerSideEncryption: aws.String("AES256"),
// 	})
// 	return fileName, err
// }

func UploudFile(buffer []byte) (string, error) {
	fileName := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	godotenv.Load(".env")

	s3Config := &aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_USER"), os.Getenv("AWS_KEY"), ""),
	}
	s3Session := session.New(s3Config)

	uploader := s3manager.NewUploader(s3Session)

	input := &s3manager.UploadInput{
		Bucket:      aws.String("bengcallbucket"),               // bucket's name
		Key:         aws.String("profile/" + fileName),          // files destination location
		Body:        bytes.NewReader(buffer),                    // content of the file
		ContentType: aws.String(http.DetectContentType(buffer)), // content type

	}
	res, err := uploader.UploadWithContext(context.Background(), input)

	return res.Location, err

}

// ParseTemplate HTML from file
func (r *RequestPDF) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.Bytes()
	return nil
}

// func newReport(name string, noInv string) *gofpdf.Fpdf {
// 	//set page invoice
// 	pdf := gofpdf.New("P", "mm", "A4", "")
// 	pdf.AddPage()
// 	pdf.SetFont("Times", "B", 16)
// 	pdf.Cell(16, 5, "Invoice Bengcall")
// 	pdf.Ln(8)

// 	pdf.SetFont("Times", "", 16)
// 	pdf.Cell(40, 10, "Time :")
// 	pdf.SetFont("Times", "", 16)
// 	pdf.Cell(40, 10, tgl)
// 	pdf.Ln(8)

// 	pdf.SetFont("Times", "", 16)
// 	pdf.Cell(40, 10, "No Invoice :")
// 	pdf.SetFont("Times", "", 16)
// 	pdf.Cell(40, 10, noInv) //no Invoice
// 	pdf.Ln(8)

// 	pdf.SetFont("Times", "", 16)
// 	pdf.Cell(40, 10, "Penerima :")
// 	pdf.SetFont("Times", "", 16)
// 	pdf.Cell(40, 10, name) //name penerima
// 	pdf.Ln(20)

// 	return pdf
// }

// func header(pdf *gofpdf.Fpdf, hdr []string) *gofpdf.Fpdf {
// 	pdf.SetFont("Times", "B", 16)
// 	pdf.SetFillColor(240, 240, 240)
// 	for _, str := range hdr {
// 		pdf.CellFormat(30, 7, str, "1", 0, "", true, 0, "")
// 	}
// 	pdf.Ln(-1)
// 	return pdf
// }

// func table(name string, pdf *gofpdf.Fpdf, tbl [][]string) *gofpdf.Fpdf {
// 	pdf.SetFont("Times", "", 16)
// 	pdf.SetFillColor(240, 240, 240)

// 	// Every column gets aligned according to its contents.
// 	align := []string{"C", "C", "L"}
// 	for _, line := range tbl {
// 		for i, str := range line {
// 			pdf.CellFormat(30, 7, str, "B", 0, align[i], false, 0, "")
// 		}
// 		pdf.Ln(-1)
// 	}
// 	//text dibawah tabel
// 	pdf.SetFont("Times", "", 16)
// 	pdf.Cell(40, 10, name)
// 	return pdf
// }

// func image(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
// 	pdf.ImageOptions("merdeka.png", 170, 10, 25, 25, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
// 	return pdf
// }

// func savePDF(pdf *gofpdf.Fpdf) (*gofpdf.Fpdf, error) {
// 	return pdf, pdf.OutputFileAndClose("invoice.pdf")
// }
