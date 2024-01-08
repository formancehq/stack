<?php

/**
 * Code generated by Speakeasy (https://speakeasyapi.dev). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Operations;


class V2ReadTriggerResponse
{
    /**
     * HTTP response content type for this operation
     * 
     * @var string $contentType
     */
	
    public string $contentType;
    
    /**
     * HTTP response status code for this operation
     * 
     * @var int $statusCode
     */
	
    public int $statusCode;
    
    /**
     * Raw HTTP response; suitable for custom response parsing
     * 
     * @var ?\Psr\Http\Message\ResponseInterface $rawResponse
     */
	
    public ?\Psr\Http\Message\ResponseInterface $rawResponse;
    
    /**
     * General error
     * 
     * @var ?\formance\stack\Models\Shared\V2Error $v2Error
     */
	
    public ?\formance\stack\Models\Shared\V2Error $v2Error = null;
    
    /**
     * A specific trigger
     * 
     * @var ?\formance\stack\Models\Shared\V2ReadTriggerResponse $v2ReadTriggerResponse
     */
	
    public ?\formance\stack\Models\Shared\V2ReadTriggerResponse $v2ReadTriggerResponse = null;
    
	public function __construct()
	{
		$this->contentType = "";
		$this->statusCode = 0;
		$this->rawResponse = null;
		$this->v2Error = null;
		$this->v2ReadTriggerResponse = null;
	}
}
