package adopt

// sedum:anchor:imports
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/verify.go","path":"fmt"},"kwargs":"sha256:9c4c1464fd2d77f5f1c585f313bc37fc100a0a4d35ab4b1d880d33e48156b9c7","body":"sha256:bcbac51883cbf8064e9196b371cb880cee429866652447ad557533d817857609"}
import "fmt"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/verify.go","path":"strings"},"kwargs":"sha256:8294cddd1c21ff442c65ba6e781f65a11ac3fdef782d5259c294b6c1d23aac68","body":"sha256:a9af3c0bc4bf9342df9d7e89a324523fdba7941a9ca4190270e68ff3c4a25ea6"}
import "strings"

// /sedum:addImport
// sedum:addImport {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:d6a9443c14d82b7bdf19bb96637f957fd44ad9bfbe8f167c97e6f7699c8b830f","id":{"file":"internal/adopt/verify.go","path":"github.com/livecodelife/sedum/internal/genpkg"},"kwargs":"sha256:541a490a4ad1d772f889be260e4b11fe99875dc461ea969436cbc3788c234495","body":"sha256:5b4c7dbee430534527d893269527d915be58accc7a8330a8c104c2014153525d"}
import "github.com/livecodelife/sedum/internal/genpkg"

// /sedum:addImport

// sedum:anchor:decls
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/verify.go","name":"StripMarkers"},"kwargs":"sha256:2fbf7bae6747d4c8abb06c1a5a4412939e16c573488de35559a6a9f5c75b2efb","body":"sha256:cc0dfa4d91869816d680d2aeb079be4f595ad0b47f2ab4077a00eadc43ac37d4"}
// StripMarkers removes every Sedum marker line from content, leaving all other lines unchanged.
func StripMarkers(comment genpkg.Comment, content string) string {
	if comment.Prefix == "" {
		return content
	}
	open := comment.Prefix + " sedum:"
	closing := comment.Prefix + " /sedum:"
	var b strings.Builder
	for _, line := range strings.SplitAfter(content, "\n") {
		trimmed := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmed, open) || strings.HasPrefix(trimmed, closing) {
			continue
		}
		b.WriteString(line)
	}
	return b.String()
}

// /sedum:setFunction
// sedum:setType {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:ebb283283a745fdc09f0294b6e79bd879a830829cb964550ec6b4432fee5bb72","id":{"file":"internal/adopt/verify.go","name":"VerifyError"},"kwargs":"sha256:a3d6822b928e6efb098354adb9e1c3922e85a5aa3129712e4730446cf0e1f864","body":"sha256:630ca09ed62f8a8d0ebb647a784cc34b3d72291e07986602e6b019f50657dbc3"}
// VerifyError reports the first line where a round trip differs.
type VerifyError struct {
	Line int
	Want string
	Got  string
}

// /sedum:setType
// sedum:setMethod {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:8ea07fc52bbc260836983438ef08ad02746b86a1028c8b1a51b60eeb7b83cc2b","id":{"file":"internal/adopt/verify.go","name":"Error","receiver":"VerifyError"},"kwargs":"sha256:dab4b3e1f9d7d5e49fb014fbea691be9db6e3705d17317a27c75f34df5bbad27","body":"sha256:5544c35f4e65bff86c62adaf545f04c7e4889c38ea60b6f0b79b07638d04c68a"}
// Error describes the differing line.
func (e *VerifyError) Error() string {
	return fmt.Sprintf("round trip: line %d differs: want \"%s\" got \"%s\"", e.Line, e.Want, e.Got)
}

// /sedum:setMethod
// sedum:setFunction {"v":2,"tier":"owned","record":"prov-2026-a15d039a","gen":"sha256:137e61eead5f9d18cba8b4328c62bfcd17878220afcdb70d2c4b12446a27938e","pkg":"sha256:6589611adfb6618a7920e43f951a79901428192ea95c1084fde9f222aecd9537","id":{"file":"internal/adopt/verify.go","name":"Verify"},"kwargs":"sha256:dead2e668252edd70c04e1da46c0f0483f4cb0f2ee58072dd36af8dddce3b5e8","body":"sha256:8ac94e145a2a6405978e4ccda467d1adfaf7ea36a01191d03c77f168b1f91e4c"}
// Verify checks that adopted, with its Sedum markers removed, matches original line by line.
// Blank lines are ignored and trailing whitespace is trimmed.
func Verify(comment genpkg.Comment, original, adopted string) error {
	significant := func(s string) []string {
		var out []string
		for _, line := range strings.Split(s, "\n") {
			line = strings.TrimRight(line, " \t\r")
			if line != "" {
				out = append(out, line)
			}
		}
		return out
	}
	want := significant(original)
	got := significant(StripMarkers(comment, adopted))
	for i := 0; i < len(want) || i < len(got); i++ {
		w := "<end of file>"
		if i < len(want) {
			w = want[i]
		}
		g := "<end of file>"
		if i < len(got) {
			g = got[i]
		}
		if w != g {
			return &VerifyError{Line: i + 1, Want: w, Got: g}
		}
	}
	return nil
}

// /sedum:setFunction
