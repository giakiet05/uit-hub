import re

with open('apps/fake-uit-server/internal/service/usecase.go', 'r') as f:
    content = f.read()

# Find all methods on *Service
methods = re.findall(r'func \(s \*Service\) (\w+)\((.*?)\) \(dto\.SuccessResponse, \*apperror\.AppError\)', content)

controller_methods = []
route_registrations = []

for method_name, args_str in methods:
    args = [arg.strip() for arg in args_str.split(',') if arg.strip()]
    
    # Analyze arguments
    bind_code = ""
    call_args = []
    has_student_id = False
    
    for arg in args:
        parts = arg.split()
        arg_name = parts[0]
        arg_type = parts[-1]
        
        if arg_name == 'studentID' or arg_name == 'id':
            bind_code += f'''
    studentID, ok := ctx.Get(middleware.StudentIDKey)
    if !ok {{
        dto.SendError(ctx, apperror.ErrUnauthorized)
        return
    }}
'''
            call_args.append("studentID.(string)")
            has_student_id = True
        elif arg_name == 'courseID':
            bind_code += f'''
    courseID := ctx.Param("courseId")
'''
            call_args.append("courseID")
        elif arg_name == 'assignmentID':
            bind_code += f'''
    assignmentID := ctx.Param("assignmentId")
'''
            call_args.append("assignmentID")
        elif arg_type.startswith('dto.Query') or arg_type == 'dto.QueryYearSemester' or arg_type == 'dto.QueryRoomsAvailability':
            bind_code += f'''
    var query {arg_type}
    if err := ctx.ShouldBindQuery(&query); err != nil {{
        dto.SendError(ctx, apperror.ErrBadRequest)
        return
    }}
'''
            call_args.append("query")
        elif arg_type.startswith('dto.') and 'Request' in arg_type:
            bind_code += f'''
    var req {arg_type}
    if err := ctx.ShouldBindJSON(&req); err != nil {{
        dto.SendError(ctx, apperror.ErrBadRequest)
        return
    }}
'''
            call_args.append("req")
        else:
            # Fallback
            call_args.append(f'ctx.Param("{arg_name}")')

    call_args_str = ", ".join(call_args)
    
    controller_method = f'''
func (c *UitController) {method_name}(ctx *gin.Context) {{
{bind_code}
    res, appErr := c.service.{method_name}({call_args_str})
    if appErr != nil {{
        dto.SendError(ctx, *appErr)
        return
    }}
    dto.SendSuccess(ctx, http.StatusOK, res.Data)
}}
'''
    controller_methods.append(controller_method)
    
    # Attempt to guess route based on method name
    http_method = "POST"
    if method_name.startswith("Get") or method_name.startswith("List"):
        http_method = "GET"
        
    path = f"/{method_name.lower().replace('get', '').replace('create', '').replace('submit', '')}"
    
    route_registrations.append(f'    rg.{http_method}("{path}", c.{method_name})')

controller_code = f'''package controller

import (
	"net/http"

	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/middleware"
	"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/service"
	"github.com/gin-gonic/gin"
)

type UitController struct {{
	service *service.Service
}}

func NewUitController(service *service.Service) *UitController {{
	return &UitController{{service: service}}
}}

{"".join(controller_methods)}
'''

with open('apps/fake-uit-server/internal/controller/uit_controller.go', 'w') as f:
    f.write(controller_code)

print("Generated uit_controller.go")
