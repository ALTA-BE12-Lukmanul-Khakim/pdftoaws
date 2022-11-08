package modify2

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
	"github.com/jung-kurt/gofpdf"
)

var tgl = time.Now().Format("01-02-2006")
var name = "Lukman"
var noInv = "123"
var uploader *s3manager.Uploader

func Create() {
	pdf := newReport(name, noInv)

	pdf = header(pdf, []string{"1st column", "2nd", "3rd"})

	//pdf table content
	data := [][]string{
		{"1 1", "1 2", "1 3"},
		{"2 1", "2 2", "2 3"},
		{"abc", "def", "geh"},
	}

	pdf = table(name, pdf, data)
	pdf = image(pdf)

	if pdf.Err() {
		log.Fatalf("Failed creating PDF report: %s\n", pdf.Error())
		return
	}

	err := pdf.OutputFileAndClose("report-example.pdf")
	// if err != nil {
	// 	log.Fatalf("error saving pdf file: %s", err)
	// 	return

	// }

	uploader = NewUplouder()

	res, err := UploudFile("report-example.pdf")
	if err != nil {
		log.Println("in create", err)
		return
	}

	// file, err := ioutil.ReadFile("report-example.pdf")
	// if err != nil {
	// 	log.Fatalf("error read pdf file: %s", err)
	// }

	log.Println(res)

}

func newReport(name string, noInv string) *gofpdf.Fpdf {
	//set page invoice
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Times", "B", 16)
	pdf.Cell(16, 5, "Invoice Bengcall")
	pdf.Ln(8)

	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "Time :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, tgl)
	pdf.Ln(8)

	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "No Invoice :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, noInv) //no Invoice
	pdf.Ln(8)

	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, "Penerima :")
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, name) //name penerima
	pdf.Ln(20)

	return pdf
}

func header(pdf *gofpdf.Fpdf, hdr []string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 16)
	pdf.SetFillColor(240, 240, 240)
	for _, str := range hdr {
		pdf.CellFormat(30, 7, str, "1", 0, "", true, 0, "")
	}
	pdf.Ln(-1)
	return pdf
}

func table(name string, pdf *gofpdf.Fpdf, tbl [][]string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 16)
	pdf.SetFillColor(240, 240, 240)

	// Every column gets aligned according to its contents.
	align := []string{"C", "C", "L"}
	for _, line := range tbl {
		for i, str := range line {
			pdf.CellFormat(30, 7, str, "B", 0, align[i], false, 0, "")
		}
		pdf.Ln(-1)
	}
	//text dibawah tabel
	pdf.SetFont("Times", "", 16)
	pdf.Cell(40, 10, name)
	return pdf
}

func image(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.ImageOptions("merdeka.png", 170, 10, 25, 25, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}

func NewUplouder() *s3manager.Uploader {
	godotenv.Load(".env")

	s3Config := &aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_USER"), os.Getenv("AWS_KEY"), ""),
	}
	s3Session := session.New(s3Config)
	uploader := s3manager.NewUploader(s3Session)

	return uploader
}

func UploudFile(filename string) (string, error) {
	//fileName := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	file, err := os.Open(filename)
	if err != nil {
		log.Println("in uploud", err)
		return "", err
	}
	defer file.Close()

	input := &s3manager.UploadInput{
		Bucket:             aws.String("bengcallbucket"),      // bucket's name
		Key:                aws.String("profile/" + filename), // files destination location
		Body:               file,                              // content of the file
		ContentDisposition: aws.String("attachment"),          // content type

	}
	res, err := uploader.UploadWithContext(context.Background(), input)

	return res.Location, err

}
