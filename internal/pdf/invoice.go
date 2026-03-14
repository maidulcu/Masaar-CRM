package pdf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type InvoiceData struct {
	InvoiceNo   string
	IssuedAt    time.Time
	Subtotal    float64
	VATRate     float64
	VATAmount   float64
	Total       float64
	DealTitle   string
	CompanyName string
	CompanyAddr string
	CompanyVAT  string
}

func GenerateInvoice(data InvoiceData) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Helvetica", "", 10)

	pdf.AddPage()

	pdf.SetFont("Helvetica", "B", 20)
	pdf.CellFormat(190, 10, "TAX INVOICE", "", 0, "C", false, 0, "")
	pdf.Ln(15)

	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(190, 5, data.CompanyName, "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(190, 5, data.CompanyAddr, "", 0, "L", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(190, 5, fmt.Sprintf("VAT Number: %s", data.CompanyVAT), "", 0, "L", false, 0, "")
	pdf.Ln(15)

	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(40, 7, "Invoice No:", "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(150, 7, data.InvoiceNo, "", 0, "L", false, 0, "")
	pdf.Ln(7)

	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(40, 7, "Date:", "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(150, 7, data.IssuedAt.Format("02 Jan 2006"), "", 0, "L", false, 0, "")
	pdf.Ln(7)

	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(40, 7, "Description:", "", 0, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(150, 7, data.DealTitle, "", 0, "L", false, 0, "")
	pdf.Ln(15)

	pdf.SetFillColor(245, 245, 245)
	pdf.Rect(10, pdf.GetY(), 190, 8, "F")
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(120, 8, "Description", "", 0, "L", false, 0, "")
	pdf.CellFormat(70, 8, "Amount (AED)", "", 0, "R", false, 0, "")
	pdf.Ln(8)

	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(120, 8, "Services", "", 0, "L", false, 0, "")
	pdf.CellFormat(70, 8, fmt.Sprintf("%.2f", data.Subtotal), "", 0, "R", false, 0, "")
	pdf.Ln(15)

	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(120, 6, "Subtotal:", "", 0, "L", false, 0, "")
	pdf.CellFormat(70, 6, fmt.Sprintf("%.2f", data.Subtotal), "", 0, "R", false, 0, "")
	pdf.Ln(6)

	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(120, 6, fmt.Sprintf("VAT (%.0f%%):", data.VATRate*100), "", 0, "L", false, 0, "")
	pdf.CellFormat(70, 6, fmt.Sprintf("%.2f", data.VATAmount), "", 0, "R", false, 0, "")
	pdf.Ln(8)

	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(10, pdf.GetY(), 160, pdf.GetY())
	pdf.Ln(2)

	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(120, 8, "Total (AED):", "", 0, "L", false, 0, "")
	pdf.CellFormat(70, 8, fmt.Sprintf("%.2f", data.Total), "", 0, "R", false, 0, "")
	pdf.Ln(20)

	pdf.SetFont("Helvetica", "", 8)
	pdf.CellFormat(190, 5, "This is a computer-generated invoice and does not require a signature.", "", 0, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Helvetica", "I", 8)
	pdf.CellFormat(190, 5, fmt.Sprintf("Generated on %s", time.Now().Format("02 Jan 2006 15:04")), "", 0, "C", false, 0, "")

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
