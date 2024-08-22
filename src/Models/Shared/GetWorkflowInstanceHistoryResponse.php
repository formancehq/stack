<?php

/**
 * Code generated by Speakeasy (https://speakeasy.com). DO NOT EDIT.
 */

declare(strict_types=1);

namespace formance\stack\Models\Shared;


class GetWorkflowInstanceHistoryResponse
{
    /**
     * $data
     *
     * @var array<WorkflowInstanceHistory> $data
     */
    #[\JMS\Serializer\Annotation\SerializedName('data')]
    #[\JMS\Serializer\Annotation\Type('array<\formance\stack\Models\Shared\WorkflowInstanceHistory>')]
    public array $data;

    /**
     * @param  ?array<WorkflowInstanceHistory>  $data
     */
    public function __construct(?array $data = null)
    {
        $this->data = $data;
    }
}