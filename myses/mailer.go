package myses

import (
    "bytes"
    "encoding/base64"
    "fmt"
    "html/template"
    "log"
    "mime/multipart"
    "net/textproto"
    "os"
    "regexp"
    "strings"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/ses"
)

const (
    FROM     = "codex.cryptum@acmvit.in"
    FROMNAME = "Codex Cryptum | ACM-VIT"
    SUBJECT  = "WELCOME TO CODEX CRYPTUM 3.0"
)

const emailTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>Codex Cryptum 3.0</title>
    <link rel="preconnect" href="https://fonts.googleapis.com" />
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
    <link
        href="https://fonts.googleapis.com/css2?family=Bricolage+Grotesque:opsz,wght@12..96,200..800&family=DM+Sans:ital,opsz,wght@0,9..40,100..1000;1,9..40,100..1000&family=Familjen+Grotesk:ital,wght@0,400..700;1,400..700&display=swap"
        rel="stylesheet" />
</head>
<body
    style="text-align: center; font-family: 'Bricolage Grotesque', sans-serif; background-color: #000000; color: #ff7c00; margin: 0; padding: 0;">
    <!-- Top Image (Logo) -->
    <img src="cid:header.png" alt="Codex Cryptum Logo"
        style="width: 100%; max-width: 600px; margin-left: auto; margin-right: auto;" />
    <!-- Main Content -->
    <div style="padding: 20px; text-align: left; max-width: 600px; margin: 0 auto;">
        <!-- Title and Greeting -->
        <p style="font-size: 16px; color: #ff7c00;">Greetings from ACM-VIT!</p>
        <p style="font-size: 16px; color: #ff7c00;">Dear {{.ParticipantName}},</p>
        <p style="font-size: 16px; color: #ff7c00;">
            Thank you for registering for <strong style="color: #ff7c00;">Codex Cryptum 3.0!</strong> We're excited to
            have you join us for an insightful session on <span style="color: #ff7c00;">cybersecurity and
                cryptography</span>.
        </p>
        <p style="font-size: 16px; color: #ff7c00;">
            The event will begin with a talk by <strong style="color: #ff7c00;">Mr. Subhash Choudhary</strong>, CTO of
            <strong style="color: #ff7c00;">Dukaan</strong>, where he’ll share his journey from a technical lead to
            co-founding one of India’s leading e-commerce platforms. Following this, you'll dive into a hands-on
            workshop on <span style="color: #ff7c00;">cybersecurity and cryptography</span>, exploring vital skills and
            insights in these fields.
        </p>
        <!-- Event Particulars -->
        <h2 style="color: #ff7c00; font-size: 24px; margin-bottom: 10px;">Event Particulars:</h2>
        <ul style="font-size: 16px; color: #ff7c00; list-style-type: none; padding: 0;">
            <li><strong style="color: #ff7c00;">Date:</strong> 22 September, 2024</li>
            <li><strong style="color: #ff7c00;">Timings:</strong> 9 AM - 6 PM</li>
            <li><strong style="color: #ff7c00;">Venue:</strong> Kamaraj Auditorium, 7th Floor, Technology Tower</li>
        </ul>
        <!-- Reminder -->
        <p style="font-size: 16px; color: #ff7c00;"><strong>Don’t forget to bring your laptops and ID cards.</strong>
        </p>
        <p style="font-size: 16px; color: #ff7c00;">
            Be on time if you don’t want to miss an amazing ice-breaking session we have planned for you all.
        </p>
        <!-- Social Media and Final Remarks -->
        <p style="font-size: 16px; color: #ff7c00;">
            For further updates, follow us on Instagram at <strong><a style="color: #ff7c00;"
                    href="https://www.instagram.com/acmvit/" target="_blank">@acmvit</a></strong>.
        </p>
        <p style="font-size: 16px; color: #ff7c00;">
            If you have any questions or need further information, feel free to reach out. We look forward to seeing you
            there!
        </p>
        <a href="https://www.instagram.com/acmvit/" target="_blank">
            <div
                style="font-size: 30px; text-align: center; color: #ff7c00; font-weight: bold; border-radius: 10px; background-color: #F3E8B5; vertical-align:middle;">
                Follow us on Instagram
            </div>
        </a>
        <!-- Sign-Off -->
        <p style="font-size: 16px; color: #ff7c00; font-weight: bold;">
            Thanks and Best Regards,<br />Manav Muthanna<br />Chairperson<br />
            <strong>Association for Computing Machinery, VIT</strong>
        </p>
    </div>
    <!-- Bottom Image -->
    <img src="cid:footer.png" alt="Gravitas 2024 Logo" style="width: 100%; max-width: 600px; margin:20px;" />
</body>
</html>
`

func sendEmail(sesClient *ses.SES, email string, data map[string]interface{}) error {
    var htmlBody bytes.Buffer
    writer := multipart.NewWriter(&htmlBody)

    htmlBody.WriteString("From: " + FROM + "\r\n")
    htmlBody.WriteString("To: " + email + "\r\n")
    htmlBody.WriteString("Subject: " + SUBJECT + "\r\n")
    htmlBody.WriteString("MIME-Version: 1.0\r\n")
    htmlBody.WriteString("Content-Type: multipart/related; boundary=" + writer.Boundary() + "\r\n")
    htmlBody.WriteString("\r\n") // End of headers

    // Write the HTML part
    htmlPartHeaders := textproto.MIMEHeader{}
    htmlPartHeaders.Set("Content-Type", "text/html; charset=UTF-8")
    htmlPart, err := writer.CreatePart(htmlPartHeaders)
    if err != nil {
        return err
    }

    tmpl, err := template.New("email").Parse(emailTemplate)
    if err != nil {
        return err
    }

    err = tmpl.Execute(htmlPart, data)
    if err != nil {
        return err
    }

    // Write the inline images
    images := []struct {
        Path     string
        CID      string
        MimeType string
    }{
        {"myses/footer.png", "footer.png", "image/png"},
        {"myses/header.png", "header.png", "image/png"},
    }

    for _, img := range images {
        imageBytes, err := os.ReadFile(img.Path)
        if err != nil {
            return err
        }
        imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)

        imagePartHeaders := textproto.MIMEHeader{}
        imagePartHeaders.Set("Content-Type", img.MimeType)
        imagePartHeaders.Set("Content-Transfer-Encoding", "base64")
        imagePartHeaders.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", img.Path))
        imagePartHeaders.Set("Content-ID", fmt.Sprintf("<%s>", img.CID))

        imagePart, err := writer.CreatePart(imagePartHeaders)
        if err != nil {
            return err
        }
        imagePart.Write([]byte(imageBase64))
    }

    err = writer.Close()
    if err != nil {
        return err
    }

    // Trim any whitespace or control characters from FROM and FROMNAME
    trimmedFrom := strings.TrimSpace(FROM)
    trimmedFromName := strings.TrimSpace(FROMNAME)

    source := fmt.Sprintf("%s <%s>", trimmedFromName, trimmedFrom)
    log.Printf("Source email: %s", source) // Debug log

    // Validate email format
    emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
    if !regexp.MustCompile(emailRegex).MatchString(trimmedFrom) {
        log.Fatalf("Invalid email format: %s", trimmedFrom)
    }

    input := &ses.SendEmailInput{
        Destination: &ses.Destination{
            ToAddresses: []*string{aws.String(email)},
        },
        Message: &ses.Message{
            Body: &ses.Body{
                Html: &ses.Content{
                    Data: aws.String(htmlBody.String()),
                },
            },
            Subject: &ses.Content{
                Data: aws.String(SUBJECT),
            },
        },
        Source: aws.String(source),
    }

    _, err = sesClient.SendEmail(input)
    return err
}

func SendSESEmail(email string, sesClient *ses.SES, data map[string]interface{}) {
    err := sendEmail(sesClient, email, data)
    if err != nil {
        log.Fatalf("Failed to send email: %v", err)
    }
    log.Println("Email sent!")
}